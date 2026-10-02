package audit

import "context"

type aiInvocationKey struct{}
type aiInvocation struct{ runID, toolID string }

func WithAIInvocation(ctx context.Context, runID, toolID string) context.Context {
	if !safeAIIdentifier(runID) || !safeAIIdentifier(toolID) {
		return ctx
	}
	return context.WithValue(ctx, aiInvocationKey{}, aiInvocation{runID: runID, toolID: toolID})
}

func withAIContext(ctx context.Context, metadata Metadata) Metadata {
	metadata.RunID, metadata.ToolID = "", ""
	invocation, ok := ctx.Value(aiInvocationKey{}).(aiInvocation)
	if ok {
		metadata.RunID, metadata.ToolID = invocation.runID, invocation.toolID
	}
	return metadata
}

func safeAIIdentifier(value string) bool {
	if value == "" || len(value) > 128 {
		return false
	}
	for _, char := range value {
		if !safeAICharacter(char) {
			return false
		}
	}
	return true
}

func safeAICharacter(char rune) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' ||
		char >= '0' && char <= '9' || char == '_' || char == '-' || char == '.'
}
