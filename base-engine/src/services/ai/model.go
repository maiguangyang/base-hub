package ai

import (
	"context"
	"net/http"

	"base-engine/agentkit/chatmodel"
	"google.golang.org/adk/v2/model"
)

type ModelConfig = chatmodel.ModelConfig

var (
	ErrModelConfig      = chatmodel.ErrModelConfig
	ErrModelProtocol    = chatmodel.ErrModelProtocol
	ErrModelUnavailable = chatmodel.ErrModelUnavailable
)

func NewModel(ctx context.Context, cfg ModelConfig) (model.LLM, error) {
	return chatmodel.NewModel(ctx, cfg)
}

func newModelWithClient(ctx context.Context, cfg ModelConfig, supplied *http.Client) (model.LLM, error) {
	return chatmodel.NewModelWithClient(ctx, cfg, supplied)
}
