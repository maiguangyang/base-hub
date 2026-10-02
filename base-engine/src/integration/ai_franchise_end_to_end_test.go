package integration_test

import (
	"context"
	"encoding/json"
	"iter"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"base-engine/gen"
	aitools "base-engine/src/services/ai/tools"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type scriptedFranchiseModel struct{ calls int }

func (*scriptedFranchiseModel) Name() string { return "scripted-franchise-chat-completions" }

func (script *scriptedFranchiseModel) GenerateContent(_ context.Context, _ *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		script.calls++
		var part *genai.Part
		switch script.calls {
		case 1:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "select-read", Name: "select_tools", Args: map[string]any{"names": []string{"FranchiseStores"}}}}
		case 2:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "read-1", Name: "FranchiseStores", Args: map[string]any{"page": 1, "pageSize": 10}}}
		case 3:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "plan-1", Name: "propose_plan", Args: map[string]any{"steps": []any{map[string]any{"operationId": "graphql.mutation.updateStore", "toolId": "FranchiseUpdateStore", "arguments": map[string]any{"id": "store-a-draft", "input": map[string]any{"name": "Updated by AI"}}, "maxCalls": 1, "sequence": 1}}}}}
		case 4:
			part = genai.NewPartFromText("Preview ready")
		case 5:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "select-write", Name: "select_tools", Args: map[string]any{"names": []string{"FranchiseUpdateStore"}}}}
		case 6:
			part = &genai.Part{FunctionCall: &genai.FunctionCall{ID: "write-1", Name: "FranchiseUpdateStore", Args: map[string]any{"id": "store-a-draft", "input": map[string]any{"name": "Updated by AI"}}}}
		default:
			part = genai.NewPartFromText("Finished")
		}
		yield(&model.LLMResponse{Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{part}}, TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 20, CandidatesTokenCount: 3, TotalTokenCount: 23}}, nil)
	}
}

func TestAIFranchiseEndToEndReadApprovalWriteAndAudit(t *testing.T) {
	fixture := newSecurityFixture(t)
	fixture.createAll([]gen.Store{{ID: "store-a-draft", Code: "A-DRAFT", Name: "Draft", Lifecycle: gen.StoreLifecycleDraft, OrganizationID: "org-a"}})
	router, script := franchiseAIRouter(t, fixture)
	preview := callAIIntegrationWithSession(t, router, fixture, "session-a-1", "/api/ai/preview", `{"prompt":"查看门店并将 store-a-draft 改名为 Updated by AI"}`)
	if preview.Code != http.StatusOK || !strings.Contains(preview.Body.String(), `"toolId":"FranchiseStores"`) || !strings.Contains(preview.Body.String(), `"toolId":"FranchiseUpdateStore"`) {
		t.Fatalf("franchise preview = %d %s", preview.Code, preview.Body.String())
	}
	token := previewTokenFromSSE(t, preview.Body.String())
	request, _ := json.Marshal(map[string]string{"previewToken": token})
	run := callAIIntegrationWithSession(t, router, fixture, "session-a-1", "/api/ai/run", string(request))
	if run.Code != http.StatusOK || !strings.Contains(run.Body.String(), `"status":"SUCCESS"`) || script.calls != 7 {
		t.Fatalf("franchise run = %d %s, model calls = %d", run.Code, run.Body.String(), script.calls)
	}
	if again := callAIIntegrationWithSession(t, router, fixture, "session-a-1", "/api/ai/run", string(request)); again.Code == http.StatusOK {
		t.Fatal("replayed franchise plan token accepted")
	}
	assertAIFranchiseStoreAndAudit(t, fixture)
}

func franchiseAIRouter(t *testing.T, fixture *securityFixture) (*mux.Router, *scriptedFranchiseModel) {
	t.Helper()
	script := &scriptedFranchiseModel{}
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
	return router, script
}

func assertAIFranchiseStoreAndAudit(t *testing.T, fixture *securityFixture) {
	t.Helper()
	var store gen.Store
	if err := fixture.db.First(&store, "id = ?", "store-a-draft").Error; err != nil {
		t.Fatal(err)
	}
	if store.Name != "Updated by AI" {
		t.Fatalf("franchise store name = %q", store.Name)
	}
	assertAIFranchiseAudit(t, fixture)
}

func assertAIFranchiseAudit(t *testing.T, fixture *securityFixture) {
	t.Helper()
	var record gen.AuditLog
	if err := fixture.db.Where("action = ? AND resource_id = ?", "store:update", "store-a-draft").First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.ActorAccountID == nil || *record.ActorAccountID != "account-shared" || record.SessionID == nil || *record.SessionID != "session-a-1" {
		t.Fatalf("franchise audit identity = %+v", record)
	}
	if record.MetadataJSON == nil || !strings.Contains(*record.MetadataJSON, `"toolId":"FranchiseUpdateStore"`) || !strings.Contains(*record.MetadataJSON, `"runId":`) {
		t.Fatalf("franchise audit metadata = %+v", record.MetadataJSON)
	}
}
