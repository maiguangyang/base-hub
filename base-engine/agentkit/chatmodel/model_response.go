package chatmodel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

const maxModelResponseBytes = 8 * 1024 * 1024

type chatChoice struct {
	Index        *int         `json:"index"`
	Message      *chatMessage `json:"message"`
	FinishReason *string      `json:"finish_reason"`
}

type chatCompletion struct {
	Choices []chatChoice `json:"choices"`
	Usage   *chatUsage   `json:"usage"`
	Model   string       `json:"model"`
}

func readChatCompletion(reader io.Reader, declared map[string]struct{}, yield func(*model.LLMResponse, error) bool) error {
	completion, err := decodeChatCompletion(reader)
	if err != nil {
		return err
	}
	choice := completion.Choices[0]
	if choice.Index == nil || *choice.Index != 0 || choice.Message == nil ||
		choice.Message.Role != "assistant" || choice.FinishReason == nil {
		return ErrModelProtocol
	}
	parts, err := responseParts(choice.Message.Content, choice.Message.ReasoningContent, choice.Message.ToolCalls, declared, *choice.FinishReason)
	if err != nil {
		return err
	}
	response := &model.LLMResponse{
		Content:      &genai.Content{Role: genai.RoleModel, Parts: parts},
		TurnComplete: true, UsageMetadata: completion.Usage.metadata(), ModelVersion: completion.Model,
	}
	if !yield(response, nil) {
		return errStreamStopped
	}
	return nil
}

func decodeChatCompletion(reader io.Reader) (chatCompletion, error) {
	data, err := io.ReadAll(io.LimitReader(reader, maxModelResponseBytes+1))
	if err != nil || len(data) > maxModelResponseBytes || !utf8.Valid(data) {
		return chatCompletion{}, ErrModelProtocol
	}
	var completion chatCompletion
	if decodeSingleJSON(data, &completion) != nil || len(completion.Choices) != 1 || !completion.Usage.valid() {
		return chatCompletion{}, ErrModelProtocol
	}
	return completion, nil
}

func responseParts(text, reasoning string, calls []chatToolCall, declared map[string]struct{}, finish string) ([]*genai.Part, error) {
	if (len(calls) == 0 && finish != "stop") || (len(calls) > 0 && finish != "tool_calls") {
		return nil, fmt.Errorf("unsupported finish reason: %w", ErrModelProtocol)
	}
	if !hasVisibleResponse(text, calls) {
		return nil, fmt.Errorf("empty final response: %w", ErrModelProtocol)
	}
	var parts []*genai.Part
	if reasoning != "" {
		parts = append(parts, &genai.Part{Text: reasoning, Thought: true})
	}
	if text != "" {
		parts = append(parts, genai.NewPartFromText(text))
	}
	for _, call := range calls {
		part, err := functionCallPart(call, declared)
		if err != nil {
			return nil, fmt.Errorf("invalid function call: %w", err)
		}
		parts = append(parts, part)
	}
	return parts, nil
}

func hasVisibleResponse(text string, calls []chatToolCall) bool {
	return text != "" || len(calls) > 0
}

func functionCallPart(call chatToolCall, declared map[string]struct{}) (*genai.Part, error) {
	if call.ID == "" || call.Type != "function" || call.Function.Name == "" {
		return nil, fmt.Errorf("missing function call metadata: %w", ErrModelProtocol)
	}
	if _, exists := declared[call.Function.Name]; !exists {
		return nil, fmt.Errorf("function call not declared for this round: %w", ErrModelProtocol)
	}
	var args map[string]any
	if decodeSingleJSON([]byte(call.Function.Arguments), &args) != nil || args == nil {
		return nil, fmt.Errorf("invalid function call arguments: %w", ErrModelProtocol)
	}
	return &genai.Part{FunctionCall: &genai.FunctionCall{ID: call.ID, Name: call.Function.Name, Args: args}}, nil
}

func decodeSingleJSON(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return ErrModelProtocol
	}
	return nil
}

func (u *chatUsage) metadata() *genai.GenerateContentResponseUsageMetadata {
	if u == nil {
		return nil
	}
	return &genai.GenerateContentResponseUsageMetadata{
		PromptTokenCount:     u.PromptTokens,
		CandidatesTokenCount: u.CompletionTokens,
		TotalTokenCount:      u.TotalTokens,
	}
}

func (u *chatUsage) valid() bool {
	return u != nil && u.PromptTokens > 0 && u.CompletionTokens >= 0 &&
		u.TotalTokens >= u.PromptTokens+u.CompletionTokens
}
