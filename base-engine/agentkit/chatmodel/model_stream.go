package chatmodel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type chatDelta struct {
	Content          string         `json:"content"`
	ReasoningContent string         `json:"reasoning_content"`
	ToolCalls        []chatToolCall `json:"tool_calls"`
}

type chatStreamChoice struct {
	Index        *int       `json:"index"`
	Delta        *chatDelta `json:"delta"`
	FinishReason *string    `json:"finish_reason"`
}

type chatStreamChunk struct {
	Choices []chatStreamChoice `json:"choices"`
	Usage   *chatUsage         `json:"usage"`
	Model   string             `json:"model"`
}

type streamAccumulator struct {
	text      string
	reasoning string
	calls     map[int]*chatToolCall
	finish    string
	usage     *chatUsage
	modelName string
	declared  map[string]struct{}
	yield     func(*model.LLMResponse, error) bool
}

func readChatStream(ctx context.Context, reader io.Reader, declared map[string]struct{}, yield func(*model.LLMResponse, error) bool) error {
	state := streamAccumulator{calls: make(map[int]*chatToolCall), declared: declared, yield: yield}
	err := parseModelSSE(reader, func(data []byte) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		chunk, err := decodeChatStreamChunk(data)
		if err != nil {
			return err
		}
		if err := state.consume(chunk); err != nil {
			return fmt.Errorf("chat stream chunk: %w", err)
		}
		return nil
	})
	if errors.Is(err, errStreamStopped) {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return fmt.Errorf("chat stream frames: %w", err)
	}
	if err := state.complete(); err != nil {
		return fmt.Errorf("chat stream completion: %w", err)
	}
	return nil
}

func decodeChatStreamChunk(data []byte) (chatStreamChunk, error) {
	var chunk chatStreamChunk
	if err := decodeSingleJSON(data, &chunk); err != nil {
		return chatStreamChunk{}, describeChatStreamJSON(err, len(data))
	}
	return chunk, nil
}

func describeChatStreamJSON(err error, length int) error {
	var syntax *json.SyntaxError
	if errors.As(err, &syntax) {
		return fmt.Errorf("chat stream JSON syntax at byte %d/%d: %w", syntax.Offset, length, ErrModelProtocol)
	}
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return fmt.Errorf("chat stream JSON incomplete frame length %d: %w", length, ErrModelProtocol)
	}
	var fieldType *json.UnmarshalTypeError
	if errors.As(err, &fieldType) {
		return describeChatStreamFieldType(fieldType, length)
	}
	if errors.Is(err, ErrModelProtocol) {
		return fmt.Errorf("chat stream JSON multiple values in frame length %d: %w", length, ErrModelProtocol)
	}
	return fmt.Errorf("chat stream JSON frame length %d: %w", length, ErrModelProtocol)
}

func describeChatStreamFieldType(fieldType *json.UnmarshalTypeError, length int) error {
	kind := strings.Fields(fieldType.Value)
	if len(kind) == 0 {
		kind = []string{"unknown"}
	}
	field := fieldType.Field
	if field == "" {
		field = "root"
	}
	return fmt.Errorf("chat stream JSON field %s type %s at byte %d/%d: %w", field, kind[0], fieldType.Offset, length, ErrModelProtocol)
}

func (s *streamAccumulator) consume(chunk chatStreamChunk) error {
	if len(chunk.Choices) == 0 {
		return s.consumeUsage(chunk.Usage)
	}
	if chunk.Usage != nil && len(chunk.Choices) == 1 && s.finish != "" {
		return s.consumeTerminalUsage(chunk)
	}
	return s.consumeChoice(chunk)
}

func (s *streamAccumulator) consumeTerminalUsage(chunk chatStreamChunk) error {
	choice := chunk.Choices[0]
	if choice.Index == nil || *choice.Index != 0 || choice.Delta == nil ||
		choice.Delta.Content != "" || len(choice.Delta.ToolCalls) != 0 || choice.FinishReason != nil {
		return ErrModelProtocol
	}
	s.reasoning += choice.Delta.ReasoningContent
	return s.consumeUsage(chunk.Usage)
}

func (s *streamAccumulator) consumeUsage(usage *chatUsage) error {
	if !usage.valid() || s.usage != nil {
		return ErrModelProtocol
	}
	s.usage = usage
	return nil
}

func (s *streamAccumulator) consumeChoice(chunk chatStreamChunk) error {
	if len(chunk.Choices) != 1 || s.finish != "" || s.usage != nil {
		return ErrModelProtocol
	}
	choice := chunk.Choices[0]
	if choice.Index == nil || *choice.Index != 0 || choice.Delta == nil {
		return ErrModelProtocol
	}
	if chunk.Model != "" {
		s.modelName = chunk.Model
	}
	if err := s.consumeDelta(choice.Delta); err != nil {
		return err
	}
	if choice.FinishReason != nil {
		s.finish = *choice.FinishReason
	}
	return s.consumeChoiceUsage(chunk.Usage)
}

func (s *streamAccumulator) consumeChoiceUsage(usage *chatUsage) error {
	if usage == nil {
		return nil
	}
	if s.finish == "" {
		return ErrModelProtocol
	}
	return s.consumeUsage(usage)
}

func (s *streamAccumulator) consumeDelta(delta *chatDelta) error {
	if delta.Content != "" {
		s.text += delta.Content
		partial := &model.LLMResponse{
			Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{genai.NewPartFromText(delta.Content)}},
			Partial: true,
		}
		if !s.yield(partial, nil) {
			return errStreamStopped
		}
	}
	s.reasoning += delta.ReasoningContent
	for _, fragment := range delta.ToolCalls {
		if err := s.consumeToolFragment(fragment); err != nil {
			return err
		}
	}
	return nil
}

func (s *streamAccumulator) consumeToolFragment(fragment chatToolCall) error {
	if fragment.Index == nil || *fragment.Index < 0 || *fragment.Index >= 100 {
		return ErrModelProtocol
	}
	index := *fragment.Index
	call := s.calls[index]
	if call == nil {
		call = &chatToolCall{}
		s.calls[index] = call
	}
	if mergeFragmentField(&call.ID, fragment.ID) != nil ||
		mergeFragmentField(&call.Type, fragment.Type) != nil ||
		mergeFragmentField(&call.Function.Name, fragment.Function.Name) != nil {
		return ErrModelProtocol
	}
	call.Function.Arguments += fragment.Function.Arguments
	return nil
}

func mergeFragmentField(current *string, incoming string) error {
	if incoming == "" {
		return nil
	}
	if *current != "" && *current != incoming {
		return ErrModelProtocol
	}
	*current = incoming
	return nil
}

func (s *streamAccumulator) complete() error {
	if s.usage == nil {
		return fmt.Errorf("missing usage: %w", ErrModelProtocol)
	}
	if s.finish == "" {
		return fmt.Errorf("missing finish reason: %w", ErrModelProtocol)
	}
	calls := make([]chatToolCall, 0, len(s.calls))
	for index := 0; index < len(s.calls); index++ {
		call := s.calls[index]
		if call == nil {
			return fmt.Errorf("non-contiguous tool calls: %w", ErrModelProtocol)
		}
		calls = append(calls, *call)
	}
	parts, err := responseParts(s.text, s.reasoning, calls, s.declared, s.finish)
	if err != nil {
		return fmt.Errorf("invalid final response: %w", err)
	}
	response := &model.LLMResponse{
		Content:      &genai.Content{Role: genai.RoleModel, Parts: parts},
		TurnComplete: true, UsageMetadata: s.usage.metadata(), ModelVersion: s.modelName,
	}
	if !s.yield(response, nil) {
		return errStreamStopped
	}
	return nil
}
