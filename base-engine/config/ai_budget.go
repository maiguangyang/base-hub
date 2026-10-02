package config

import (
	"errors"
	"os"
	"strconv"
)

type AIBudgetConfig struct{ ModelContextTokens, InputTokenBudget, OutputTokenBudget int }

func LoadAIBudgetConfig() (AIBudgetConfig, error) {
	contextTokens, err := aiBudgetValue("AI_MODEL_CONTEXT_TOKENS", 32768)
	if err != nil {
		return AIBudgetConfig{}, err
	}
	inputTokens, err := aiBudgetValue("AI_INPUT_TOKEN_BUDGET", 64000)
	if err != nil {
		return AIBudgetConfig{}, err
	}
	outputTokens, err := aiBudgetValue("AI_OUTPUT_TOKEN_BUDGET", 12000)
	if err != nil {
		return AIBudgetConfig{}, err
	}
	return AIBudgetConfig{contextTokens, inputTokens, outputTokens}, nil
}

func aiBudgetValue(name string, fallback int) (int, error) {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 256 || value > 1_000_000 {
		return 0, errors.New("INVALID_AI_TOKEN_BUDGET")
	}
	return value, nil
}
