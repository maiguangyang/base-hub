package config

import "testing"

func TestAIBudgetDefaultsAllowMultiToolRuns(t *testing.T) {
	t.Setenv("AI_MODEL_CONTEXT_TOKENS", "")
	t.Setenv("AI_INPUT_TOKEN_BUDGET", "")
	t.Setenv("AI_OUTPUT_TOKEN_BUDGET", "")
	got, err := LoadAIBudgetConfig()
	if err != nil {
		t.Fatal(err)
	}
	if got.ModelContextTokens != 32768 || got.InputTokenBudget != 64000 || got.OutputTokenBudget != 12000 {
		t.Fatalf("default AI budgets = %+v", got)
	}
}
