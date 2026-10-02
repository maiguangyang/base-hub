package ai

import (
	"context"
	"iter"
	"net/http"
	"testing"
	"time"

	"base-engine/auth"
	"google.golang.org/adk/v2/model"
)

type serviceTestModel struct{}

func (serviceTestModel) Name() string { return "test" }
func (serviceTestModel) GenerateContent(context.Context, *model.LLMRequest, bool) iter.Seq2[*model.LLMResponse, error] {
	return func(func(*model.LLMResponse, error) bool) {}
}

func TestAIServiceRequiresSingleStartupInjection(t *testing.T) {
	service := newServiceFixture(t)
	assertServiceInjection(t, service)
}

func newServiceFixture(t *testing.T) *Service {
	t.Helper()
	catalog, err := NewCatalog(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewService(ServiceConfig{Model: serviceTestModel{}, Catalog: catalog, Prompt: "test prompt", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return nil, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func assertServiceInjection(t *testing.T, service *Service) {
	t.Helper()
	if service.Ready() {
		t.Fatal("service ready before handler/tools injection")
	}
	if err := service.SetProtectedHandler(nil); err == nil {
		t.Fatal("nil handler accepted")
	}
	if err := service.SetProtectedHandler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})); err != nil {
		t.Fatal(err)
	}
	if service.Ready() {
		t.Fatal("service ready before tools injection")
	}
	if err := service.SetTools(nil); err != nil {
		t.Fatal(err)
	}
	if !service.Ready() {
		t.Fatal("service not ready after complete injection")
	}
	if err := service.SetTools(nil); err == nil {
		t.Fatal("duplicate tools injection accepted")
	}
	if err := service.SetProtectedHandler(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})); err == nil {
		t.Fatal("duplicate handler injection accepted")
	}
}

func TestAIServiceRejectsInvalidBudgets(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	_, err := NewService(ServiceConfig{Model: serviceTestModel{}, Catalog: catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return nil, nil }, ModelContextTokens: 1000, InputTokenBudget: 0, OutputTokenBudget: 100})
	if err == nil {
		t.Fatal("zero cumulative input budget accepted")
	}
}

func TestAIServiceAllowsCumulativeRunBudgetAboveSingleCallContext(t *testing.T) {
	catalog, _ := NewCatalog(nil, nil)
	_, err := NewService(ServiceConfig{Model: serviceTestModel{}, Catalog: catalog, Prompt: "test",
		ResolvePrincipal:   func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return nil, nil },
		ModelContextTokens: 10000, InputTokenBudget: 64000, OutputTokenBudget: 12000})
	if err != nil {
		t.Fatalf("valid cumulative run budget rejected: %v", err)
	}
}

func TestAIServiceSharesAccountSlotAndLimitsEachPhase(t *testing.T) {
	service := newServiceFixture(t)
	release, err := service.acquire("account-1", PhasePreview)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.acquire("account-1", PhaseRun); err == nil {
		t.Fatal("concurrent run accepted")
	}
	release()
	for range 9 {
		release, err := service.acquire("account-1", PhasePreview)
		if err != nil {
			t.Fatal(err)
		}
		release()
	}
	if _, err := service.acquire("account-1", PhasePreview); err == nil {
		t.Fatal("eleventh preview accepted")
	}
	release, err = service.acquire("account-1", PhaseRun)
	if err != nil {
		t.Fatal("run incorrectly shared preview quota: ", err)
	}
	release()
}

func TestAIServiceRejectsIncompleteToolRegistration(t *testing.T) {
	approval, principal := approvalFixture(t)
	service, err := NewService(ServiceConfig{Model: serviceTestModel{}, Catalog: approval.catalog, Prompt: "test", ResolvePrincipal: func(context.Context, string, time.Time) (*auth.WorkspacePrincipal, error) { return principal, nil }, ModelContextTokens: 10000, InputTokenBudget: 5000, OutputTokenBudget: 2000})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.SetTools(nil); err == nil {
		t.Fatal("catalog started with missing tools")
	}
}
