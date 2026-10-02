package main

import (
	"context"
	_ "embed"
	"encoding/json"

	"base-engine/config"
	"base-engine/src"
	"base-engine/src/services/ai"
	aitools "base-engine/src/services/ai/tools"
	"base-engine/system_prompt"
	"google.golang.org/adk/v2/model"
)

//go:embed tools/contract_inventory.json
var contractInventoryJSON []byte

func newAICatalog() (*ai.Catalog, error) {
	var inventory struct {
		Operations []ai.ContractRecord `json:"operations"`
	}
	if err := json.Unmarshal(contractInventoryJSON, &inventory); err != nil {
		return nil, err
	}
	specs, err := aitools.PreparedSpecs()
	if err != nil {
		return nil, err
	}
	return ai.NewCatalog(inventory.Operations, specs)
}

func newAIService(store *ai.ModelConfigStore, dependencies src.Dependencies) (*ai.Service, error) {
	budget, err := config.LoadAIBudgetConfig()
	if err != nil {
		return nil, err
	}
	catalog, err := newAICatalog()
	if err != nil {
		return nil, err
	}
	return ai.NewService(ai.ServiceConfig{Catalog: catalog, Prompt: system_prompt.AdminAgent,
		ResolvePrincipal:        dependencies.Principal.Resolve,
		ValidateImageAttachment: dependencies.ProductCatalog.ValidateMainImageAttachment,
		ModelContextTokens:      budget.ModelContextTokens, InputTokenBudget: budget.InputTokenBudget, OutputTokenBudget: budget.OutputTokenBudget,
		ModelProvider: func(ctx context.Context) (model.LLM, uint64, error) {
			modelConfig, version, err := store.ActiveWithVersion(ctx)
			if err != nil {
				return nil, 0, err
			}
			instance, err := ai.NewModel(ctx, modelConfig)
			return instance, version, err
		}})
}
