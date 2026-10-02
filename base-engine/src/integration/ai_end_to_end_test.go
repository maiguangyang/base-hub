package integration_test

import (
	"context"
	"encoding/json"
	"iter"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/config"
	"base-engine/gen"
	enginesrc "base-engine/src"
	"base-engine/src/services/ai"
	aitools "base-engine/src/services/ai/tools"
	"base-engine/src/services/session"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type scriptedAIModel struct {
	calls             int
	runAdvertisedRead bool
}

func (*scriptedAIModel) Name() string { return "scripted-chat-completions" }

func (script *scriptedAIModel) GenerateContent(_ context.Context, request *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		script.calls++
		if script.calls == 5 {
			script.runAdvertisedRead = requestAdvertisesTool(request, "HqFranchises")
		}
		var part *genai.Part
		switch script.calls {
		case 1:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "select-read", Name: "select_tools", Args: map[string]any{"names": []string{"HqFranchises"}}}}
		case 2:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "read-1", Name: "HqFranchises", Args: map[string]any{"page": 1, "pageSize": 10}}}
		case 3:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "plan-1", Name: "propose_plan", Args: map[string]any{"steps": []any{map[string]any{"operationId": "graphql.mutation.suspendOrganization", "toolId": "HqSuspendOrganization", "arguments": map[string]any{"input": map[string]any{"organizationId": "org-a", "reasonCode": "CONTRACT_ENDED"}}, "maxCalls": 1, "sequence": 1}}}}}
		case 4:
			part = genai.NewPartFromText("Preview ready")
		case 5:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "select-write", Name: "select_tools", Args: map[string]any{"names": []string{"HqSuspendOrganization"}}}}
		case 6:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "write-1", Name: "HqSuspendOrganization", Args: map[string]any{"input": map[string]any{"organizationId": "org-a", "reasonCode": "CONTRACT_ENDED"}}}}
		default:
			part = genai.NewPartFromText("Finished")
		}
		yield(&model.LLMResponse{Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{part}}, TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 20, CandidatesTokenCount: 3, TotalTokenCount: 23}}, nil)
	}
}

func requestAdvertisesTool(request *model.LLMRequest, toolID string) bool {
	for _, item := range request.Config.Tools {
		for _, declaration := range item.FunctionDeclarations {
			if declaration.Name == "select_tools" && strings.Contains(declaration.Description, toolID) {
				return true
			}
		}
	}
	return false
}

func TestAIEndToEndOneTimeApprovalAndAudit(t *testing.T) {
	fixture := newSecurityFixture(t)
	script := &scriptedAIModel{}
	service := newAIIntegrationService(t, fixture, script, "Use tools and propose the write plan")
	router := mux.NewRouter()
	service.RegisterRoutes(router, fixture.cfg)
	if err := service.SetProtectedHandler(fixture.handler); err != nil {
		t.Fatal(err)
	}
	items, err := aitools.Build(service.FixedToolRuntime())
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SetTools(items); err != nil {
		t.Fatal(err)
	}
	preview := callAIIntegration(t, router, fixture, "/api/ai/preview", `{"prompt":"先查看加盟商，再停用 org-a，原因 CONTRACT_ENDED"}`)
	if !strings.Contains(preview.Body.String(), `"toolId":"HqFranchises"`) || !strings.Contains(preview.Body.String(), `"status":"SUCCESS"`) {
		t.Fatalf("preview = %s", preview.Body.String())
	}
	token := previewTokenFromSSE(t, preview.Body.String())
	request, _ := json.Marshal(map[string]string{"previewToken": token})
	run := callAIIntegration(t, router, fixture, "/api/ai/run", string(request))
	if !strings.Contains(run.Body.String(), `"status":"SUCCESS"`) || script.calls != 7 {
		t.Fatalf("run = %s, model calls = %d", run.Body.String(), script.calls)
	}
	if script.runAdvertisedRead {
		t.Fatal("unapproved read tool advertised during approved execution")
	}
	if again := callAIIntegration(t, router, fixture, "/api/ai/run", string(request)); again.Code == http.StatusOK {
		t.Fatal("replayed plan token accepted")
	}
	assertAISuspensionAudit(t, fixture)
}

func newAIIntegrationService(t *testing.T, fixture *securityFixture, llm model.LLM, prompt string, budgets ...config.AIBudgetConfig) *ai.Service {
	t.Helper()
	data, err := os.ReadFile("../../tools/contract_inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Operations []ai.ContractRecord `json:"operations"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	specs, err := aitools.PreparedSpecs()
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := ai.NewCatalog(inventory.Operations, specs)
	if err != nil {
		t.Fatal(err)
	}
	deps := enginesrc.NewDependencies(fixture.db, fixture.cfg, session.NewPublisher())
	budget := config.AIBudgetConfig{ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000}
	if len(budgets) > 0 {
		budget = budgets[0]
	}
	service, err := ai.NewService(ai.ServiceConfig{Model: llm, Catalog: catalog, Prompt: prompt, ResolvePrincipal: deps.Principal.Resolve,
		ModelContextTokens: budget.ModelContextTokens, InputTokenBudget: budget.InputTokenBudget, OutputTokenBudget: budget.OutputTokenBudget})
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func callAIIntegration(t *testing.T, router *mux.Router, fixture *securityFixture, path, body string) *httptest.ResponseRecorder {
	return callAIIntegrationWithSession(t, router, fixture, "session-hq", path, body)
}

func callAIIntegrationWithSession(t *testing.T, router *mux.Router, fixture *securityFixture, sessionID, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	request.Header.Set("Origin", "https://admin.example.com")
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(fixture.cookies[sessionID])
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func previewTokenFromSSE(t *testing.T, body string) string {
	t.Helper()
	for _, frame := range strings.Split(body, "\n\n") {
		if !strings.Contains(frame, "event: preview_ready") {
			continue
		}
		for _, line := range strings.Split(frame, "\n") {
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			var payload struct {
				PreviewToken string `json:"previewToken"`
			}
			if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &payload); err != nil {
				t.Fatal(err)
			}
			if payload.PreviewToken != "" {
				return payload.PreviewToken
			}
		}
	}
	t.Fatalf("preview token missing from SSE: %s", body)
	return ""
}

func assertAISuspensionAudit(t *testing.T, fixture *securityFixture) {
	t.Helper()
	var organization gen.Organization
	if err := fixture.db.First(&organization, "id = ?", "org-a").Error; err != nil {
		t.Fatal(err)
	}
	if organization.Status != gen.OrganizationStatusSuspended {
		t.Fatalf("organization status = %s", organization.Status)
	}
	assertAIActorAudit(t, fixture)
}

func assertAIActorAudit(t *testing.T, fixture *securityFixture) {
	t.Helper()
	var record gen.AuditLog
	if err := fixture.db.Where("action = ?", "organization:suspend").First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.ActorAccountID == nil || *record.ActorAccountID != "account-shared" || record.SessionID == nil || *record.SessionID != "session-hq" {
		t.Fatalf("audit identity = %+v", record)
	}
	if record.MetadataJSON == nil || !strings.Contains(*record.MetadataJSON, `"toolId":"HqSuspendOrganization"`) || !strings.Contains(*record.MetadataJSON, `"runId":`) {
		t.Fatalf("audit metadata = %+v", record.MetadataJSON)
	}
}
