package chatmodel

import (
	"encoding/json"
	"reflect"
	"strings"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type chatMessage struct {
	Role             string         `json:"role"`
	Content          string         `json:"content"`
	ReasoningContent string         `json:"reasoning_content,omitempty"`
	ToolCalls        []chatToolCall `json:"tool_calls,omitempty"`
	ToolCallID       string         `json:"tool_call_id,omitempty"`
}

type chatFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type chatToolCall struct {
	Index    *int         `json:"index,omitempty"`
	ID       string       `json:"id,omitempty"`
	Type     string       `json:"type,omitempty"`
	Function chatFunction `json:"function"`
}

type chatToolDefinition struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Parameters  any    `json:"parameters"`
	} `json:"function"`
}

type chatRequest struct {
	Model             string               `json:"model"`
	Messages          []chatMessage        `json:"messages"`
	Stream            bool                 `json:"stream"`
	StreamOptions     *chatStreamOptions   `json:"stream_options,omitempty"`
	ParallelToolCalls *bool                `json:"parallel_tool_calls,omitempty"`
	Temperature       *float32             `json:"temperature,omitempty"`
	MaxTokens         int32                `json:"max_tokens,omitempty"`
	Tools             []chatToolDefinition `json:"tools,omitempty"`
	ToolChoice        string               `json:"tool_choice,omitempty"`
}

type chatStreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type chatUsage struct {
	PromptTokens     int32 `json:"prompt_tokens"`
	CompletionTokens int32 `json:"completion_tokens"`
	TotalTokens      int32 `json:"total_tokens"`
}

func buildChatRequest(name string, req *model.LLMRequest, stream bool) (chatRequest, map[string]struct{}, error) {
	if req == nil || len(req.Contents) == 0 {
		return chatRequest{}, nil, ErrModelProtocol
	}
	wire := chatRequest{Model: name, Stream: stream}
	if stream {
		wire.StreamOptions = &chatStreamOptions{IncludeUsage: true}
	}
	if err := applyGenerationConfig(&wire, req.Config); err != nil {
		return chatRequest{}, nil, err
	}
	for _, content := range req.Contents {
		messages, err := chatMessages(content)
		if err != nil {
			return chatRequest{}, nil, err
		}
		wire.Messages = append(wire.Messages, messages...)
	}
	if len(wire.Messages) == 0 {
		return chatRequest{}, nil, ErrModelProtocol
	}
	declared, err := appendChatTools(&wire, req.Config)
	return wire, declared, err
}

func applyGenerationConfig(wire *chatRequest, cfg *genai.GenerateContentConfig) error {
	if cfg == nil {
		return nil
	}
	if cfg.ResponseJsonSchema != nil || cfg.ResponseSchema != nil ||
		(cfg.ResponseMIMEType != "" && cfg.ResponseMIMEType != "text/plain") {
		return ErrModelProtocol
	}
	wire.Temperature = cfg.Temperature
	wire.MaxTokens = cfg.MaxOutputTokens
	if cfg.SystemInstruction != nil {
		value, err := contentText(cfg.SystemInstruction)
		if err != nil {
			return err
		}
		wire.Messages = append(wire.Messages, chatMessage{Role: "system", Content: value})
	}
	return nil
}

func appendChatTools(wire *chatRequest, cfg *genai.GenerateContentConfig) (map[string]struct{}, error) {
	declared := map[string]struct{}{}
	if cfg == nil {
		return declared, nil
	}
	if cfg.ToolConfig != nil && cfg.ToolConfig.FunctionCallingConfig != nil &&
		cfg.ToolConfig.FunctionCallingConfig.Mode == genai.FunctionCallingConfigModeNone {
		wire.ToolChoice = "none"
		return declared, nil
	}
	for _, item := range cfg.Tools {
		if !functionOnlyTool(item) {
			return nil, ErrModelProtocol
		}
		for _, fn := range item.FunctionDeclarations {
			definition, err := chatToolDefinitionFor(fn, declared)
			if err != nil {
				return nil, err
			}
			wire.Tools = append(wire.Tools, definition)
		}
	}
	if len(wire.Tools) > 0 {
		noParallel := false
		wire.ParallelToolCalls = &noParallel
	}
	return declared, nil
}

func functionOnlyTool(item *genai.Tool) bool {
	if item == nil {
		return false
	}
	value := reflect.ValueOf(item).Elem()
	for index := 0; index < value.NumField(); index++ {
		if value.Type().Field(index).Name != "FunctionDeclarations" && !value.Field(index).IsZero() {
			return false
		}
	}
	return len(item.FunctionDeclarations) > 0
}

func chatToolDefinitionFor(fn *genai.FunctionDeclaration, declared map[string]struct{}) (chatToolDefinition, error) {
	if fn == nil || fn.Name == "" {
		return chatToolDefinition{}, ErrModelProtocol
	}
	if _, exists := declared[fn.Name]; exists {
		return chatToolDefinition{}, ErrModelProtocol
	}
	declared[fn.Name] = struct{}{}
	definition := chatToolDefinition{Type: "function"}
	definition.Function.Name = fn.Name
	definition.Function.Description = fn.Description
	definition.Function.Parameters = fn.ParametersJsonSchema
	if definition.Function.Parameters == nil {
		definition.Function.Parameters = fn.Parameters
	}
	if definition.Function.Parameters == nil {
		definition.Function.Parameters = map[string]any{"type": "object"}
	}
	return definition, nil
}

func chatMessages(content *genai.Content) ([]chatMessage, error) {
	if content == nil || len(content.Parts) == 0 {
		return nil, ErrModelProtocol
	}
	role := "user"
	if content.Role == genai.RoleModel {
		role = "assistant"
	} else if content.Role != genai.RoleUser {
		return nil, ErrModelProtocol
	}
	message := chatMessage{Role: role}
	var toolResponses []chatMessage
	for _, part := range content.Parts {
		if err := appendChatPart(&message, &toolResponses, part); err != nil {
			return nil, err
		}
	}
	var messages []chatMessage
	if hasChatMessageContent(message) {
		messages = append(messages, message)
	}
	messages = append(messages, toolResponses...)
	if len(messages) == 0 {
		return nil, ErrModelProtocol
	}
	return messages, nil
}

func appendChatPart(message *chatMessage, responses *[]chatMessage, part *genai.Part) error {
	if part == nil || part.InlineData != nil || part.FileData != nil ||
		part.ExecutableCode != nil || part.CodeExecutionResult != nil {
		return ErrModelProtocol
	}
	if err := appendChatText(message, part); err != nil {
		return err
	}
	return appendChatToolPart(message, responses, part)
}

func encodeFunctionCall(call *genai.FunctionCall) (chatToolCall, error) {
	if call.ID == "" || call.Name == "" {
		return chatToolCall{}, ErrModelProtocol
	}
	args, err := json.Marshal(call.Args)
	if err != nil {
		return chatToolCall{}, ErrModelProtocol
	}
	return chatToolCall{ID: call.ID, Type: "function", Function: chatFunction{Name: call.Name, Arguments: string(args)}}, nil
}

func encodeFunctionResponse(response *genai.FunctionResponse) (chatMessage, error) {
	if response.ID == "" || response.Name == "" {
		return chatMessage{}, ErrModelProtocol
	}
	value, err := json.Marshal(response.Response)
	if err != nil {
		return chatMessage{}, ErrModelProtocol
	}
	return chatMessage{Role: "tool", ToolCallID: response.ID, Content: string(value)}, nil
}

func contentText(content *genai.Content) (string, error) {
	var text strings.Builder
	for _, part := range content.Parts {
		if part == nil || part.FunctionCall != nil || part.FunctionResponse != nil ||
			part.InlineData != nil || part.FileData != nil || part.ExecutableCode != nil || part.CodeExecutionResult != nil {
			return "", ErrModelProtocol
		}
		text.WriteString(part.Text)
	}
	return text.String(), nil
}
