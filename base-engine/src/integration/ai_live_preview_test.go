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
	"time"

	"github.com/gorilla/mux"
	"base-engine/config"
	"base-engine/gen"
	"base-engine/src/services/ai"
	aitools "base-engine/src/services/ai/tools"
	"base-engine/system_prompt"
	"google.golang.org/adk/v2/model"
)

type observedLiveModel struct {
	model.LLM
	t *testing.T
}

func (m observedLiveModel) GenerateContent(ctx context.Context, request *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		encoded, _ := json.Marshal(request)
		count := 0
		if request.Config != nil {
			for _, entry := range request.Config.Tools {
				count += len(entry.FunctionDeclarations)
			}
		}
		m.t.Logf("model boundary request bytes=%d toolDefinitions=%d", len(encoded), count)
		for response, err := range m.LLM.GenerateContent(ctx, request, stream) {
			if err != nil {
				m.t.Logf("model boundary error type=%T value=%v", err, err)
			}
			if !yield(response, err) {
				return
			}
		}
	}
}

func TestAILiveProtectedBusinessRead(t *testing.T) {
	if os.Getenv("AI_LIVE_TEST_FROM_STORE") != "true" {
		t.Skip("live model gate is disabled")
	}
	modelConfig := activeLiveModelConfig(t)
	llm, err := ai.NewModel(t.Context(), modelConfig)
	if err != nil {
		t.Fatal(err)
	}
	fixture := newSecurityFixture(t)
	budget, err := config.LoadAIBudgetConfig()
	if err != nil {
		t.Fatal(err)
	}
	service := newAIIntegrationService(t, fixture, observedLiveModel{LLM: llm, t: t}, "Answer the user's headquarters count question in Chinese. Use select_tools to enable HqFranchises, HqFranchiseStores and HqDirectStores, then read their first pages with page 1 and pageSize 10. Use each list's total to report franchise and store counts. This is read-only; do not propose any write plan.", budget)
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
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/ai/preview", strings.NewReader(`{"prompt":"查看一下当前有多少家加盟商和多少家门店了"}`)).WithContext(ctx)
	request.Header.Set("Origin", "https://admin.example.com")
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(fixture.cookies["session-hq"])
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assertLiveReadResponse(t, response)
}

func TestAILiveAccountPermissionsQuestion(t *testing.T) {
	if os.Getenv("AI_LIVE_TEST_FROM_STORE") != "true" {
		t.Skip("live model gate is disabled")
	}
	modelConfig := activeLiveModelConfig(t)
	llm, err := ai.NewModel(t.Context(), modelConfig)
	if err != nil {
		t.Fatal(err)
	}
	fixture := newSecurityFixture(t)
	budget, err := config.LoadAIBudgetConfig()
	if err != nil {
		t.Fatal(err)
	}
	service := newAIIntegrationService(t, fixture, observedLiveModel{LLM: llm, t: t}, system_prompt.AdminAgent, budget)
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
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()
	request := httptest.NewRequest(http.MethodPost, "/api/ai/preview", strings.NewReader(`{"prompt":"查看一下我的账号都有什么权限"}`)).WithContext(ctx)
	request.Header.Set("Origin", "https://admin.example.com")
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(fixture.cookies["session-hq"])
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"status":"SUCCESS"`) || strings.Contains(response.Body.String(), "event: error") {
		t.Fatalf("account permissions question: HTTP %d; events=%v", response.Code, liveEventNames(response.Body.String()))
	}
}

func liveEventNames(body string) []string {
	var names []string
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "event:") {
			names = append(names, strings.TrimSpace(strings.TrimPrefix(line, "event:")))
		}
	}
	return names
}

func assertLiveReadResponse(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	body := response.Body.String()
	if response.Code != http.StatusOK || !containsSuccessfulReads(body, "HqFranchises", "HqFranchiseStores", "HqDirectStores") {
		var events []string
		var code string
		for _, line := range strings.Split(body, "\n") {
			if strings.HasPrefix(line, "event:") {
				events = append(events, strings.TrimSpace(strings.TrimPrefix(line, "event:")))
			}
			if strings.HasPrefix(line, "data:") && strings.Contains(line, `"code":`) {
				var payload struct {
					Code string `json:"code"`
				}
				_ = json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(line, "data:"))), &payload)
				code = payload.Code
			}
		}
		t.Fatalf("live business read: HTTP %d; tool/result markers missing; events=%v; code=%s; bodyBytes=%d", response.Code, events, code, len(body))
	}
}

func containsSuccessfulReads(body string, toolIDs ...string) bool {
	completed := map[string]bool{}
	finished := false
	for _, frame := range strings.Split(body, "\n\n") {
		if strings.HasPrefix(frame, "event: error") {
			return false
		}
		if strings.HasPrefix(frame, "event: run_finished") {
			finished = strings.Contains(frame, `"status":"SUCCESS"`)
		}
		if !strings.Contains(frame, "event: tool_finished") {
			continue
		}
		for _, toolID := range toolIDs {
			if strings.Contains(frame, `"toolId":"`+toolID+`"`) && strings.Contains(frame, `"status":"SUCCESS"`) {
				completed[toolID] = true
			}
		}
	}
	return finished && len(completed) == len(toolIDs)
}

func TestAILiveReadRejectsFailureAfterSuccessfulTools(t *testing.T) {
	body := "event: tool_finished\ndata: {\"toolId\":\"HqFranchises\",\"status\":\"SUCCESS\"}\n\n" +
		"event: error\ndata: {\"code\":\"AI_MODEL_PROTOCOL\"}\n\n" +
		"event: run_finished\ndata: {\"status\":\"FAILED\"}\n\n"
	if containsSuccessfulReads(body, "HqFranchises") {
		t.Fatal("successful tool call hid failed run")
	}
}

func activeLiveModelConfig(t *testing.T) ai.ModelConfig {
	t.Helper()
	security, err := config.LoadAIModelSecurityConfig()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gen.OpenDBFromEnvVars("")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	modelConfig, err := ai.NewModelConfigStore(db.Query(), security, nil).Active(t.Context())
	if err != nil {
		t.Fatalf("active model config: %v", err)
	}
	return modelConfig
}
