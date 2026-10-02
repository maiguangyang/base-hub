package chatmodel

import "google.golang.org/genai"

func hasChatMessageContent(message chatMessage) bool {
	return message.Content != "" || message.ReasoningContent != "" || len(message.ToolCalls) > 0
}

func appendChatText(message *chatMessage, part *genai.Part) error {
	if !part.Thought {
		message.Content += part.Text
		return nil
	}
	if message.Role != "assistant" || part.FunctionCall != nil || part.FunctionResponse != nil {
		return ErrModelProtocol
	}
	message.ReasoningContent += part.Text
	return nil
}

func appendChatToolPart(message *chatMessage, responses *[]chatMessage, part *genai.Part) error {
	if part.FunctionCall != nil {
		call, err := encodeFunctionCall(part.FunctionCall)
		if err != nil {
			return err
		}
		message.Role = "assistant"
		message.ToolCalls = append(message.ToolCalls, call)
	}
	if part.FunctionResponse != nil {
		response, err := encodeFunctionResponse(part.FunctionResponse)
		if err != nil {
			return err
		}
		*responses = append(*responses, response)
	}
	return nil
}
