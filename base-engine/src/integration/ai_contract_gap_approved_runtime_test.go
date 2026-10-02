package integration_test

import (
	"context"
	"iter"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"base-engine/gen"
	aitools "base-engine/src/services/ai/tools"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

type gapDeleteModel struct {
	readTool, writeTool, approvedID, callID string
	calls                                   int
}

func (*gapDeleteModel) Name() string { return "approved-invitation-delete" }

func (script *gapDeleteModel) GenerateContent(_ context.Context, _ *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		script.calls++
		part := script.part()
		yield(&model.LLMResponse{Content: &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{part}}, TurnComplete: true,
			UsageMetadata: &genai.GenerateContentResponseUsageMetadata{PromptTokenCount: 20, CandidatesTokenCount: 3, TotalTokenCount: 23}}, nil)
	}
}

func (script *gapDeleteModel) part() *genai.Part {
	switch script.calls {
	case 1:
		return gapFunctionCall("select-read", "select_tools", map[string]any{"names": []string{script.readTool}})
	case 2:
		return gapFunctionCall("read", script.readTool, map[string]any{"id": script.approvedID})
	case 3:
		return gapFunctionCall("plan", "propose_plan", map[string]any{"steps": []any{map[string]any{
			"operationId": "graphql.mutation.deleteMembershipInvitations", "toolId": script.writeTool,
			"arguments": map[string]any{"ids": []string{script.approvedID}}, "targetIds": []string{script.approvedID}, "maxCalls": 1, "sequence": 1,
		}}})
	case 4:
		return genai.NewPartFromText("Preview ready")
	case 5:
		return gapFunctionCall("select-write", "select_tools", map[string]any{"names": []string{script.writeTool}})
	case 6:
		return gapFunctionCall("delete", script.writeTool, map[string]any{"ids": []string{script.callID}})
	default:
		return genai.NewPartFromText("Finished")
	}
}

func gapFunctionCall(id, name string, args map[string]any) *genai.Part {
	return &genai.Part{FunctionCall: &genai.FunctionCall{ID: id, Name: name, Args: args}}
}

func TestContractGapDeleteToolsUseRealApprovalRuntime(t *testing.T) {
	for _, tc := range contractGapWriteCases() {
		for _, drift := range []bool{false, true} {
			name := tc.tool + "/approved"
			if drift {
				name = tc.tool + "/target-drift"
			}
			t.Run(name, func(t *testing.T) { testGapDeleteRuntime(t, tc, drift) })
		}
	}
}

func testGapDeleteRuntime(t *testing.T, tc protectedGapWriteCase, drift bool) {
	fixture := newSecurityFixture(t)
	seedContractGapTargets(fixture)
	otherID := tc.ownID + "-other"
	membershipID := "membership-hq"
	if tc.session != "session-hq" {
		membershipID = "membership-staff-a"
	}
	fixture.createAll([]gen.MembershipInvitation{{ID: otherID, MembershipID: membershipID, InvitedByAccountID: "account-shared", ExpiresAt: time.Now().Add(time.Hour)}})
	readTool := "HqMembershipInvitation"
	if tc.session != "session-hq" {
		readTool = "FranchiseMembershipInvitation"
	}
	script := &gapDeleteModel{readTool: readTool, writeTool: tc.tool, approvedID: tc.ownID, callID: tc.ownID}
	if drift {
		script.callID = otherID
	}
	router := gapDeleteServiceRouter(t, fixture, script)
	preview := callAIIntegrationWithSession(t, router, fixture, tc.session, "/api/ai/preview", `{"prompt":"Delete the selected invitation"}`)
	token := previewTokenFromSSE(t, preview.Body.String())
	run := callAIIntegrationWithSession(t, router, fixture, tc.session, "/api/ai/run", `{"previewToken":"`+token+`"}`)
	if script.calls < 6 {
		t.Fatalf("write tool was never attempted: calls=%d body=%s", script.calls, run.Body.String())
	}
	if drift {
		if !strings.Contains(run.Body.String(), `"status":"FAILED"`) {
			t.Fatalf("changed target accepted: %s", run.Body.String())
		}
		assertInvitationDeleteState(t, fixture, tc.ownID, false)
		assertInvitationDeleteState(t, fixture, otherID, false)
		assertAuditCount(t, fixture, tc.action, 0)
		return
	}
	if !strings.Contains(run.Body.String(), `"status":"SUCCESS"`) {
		t.Fatalf("approved deletion failed: %s", run.Body.String())
	}
	assertInvitationDeleteState(t, fixture, tc.ownID, true)
	assertInvitationDeleteState(t, fixture, otherID, false)
	assertAuditCount(t, fixture, tc.action, 1)
}

func gapDeleteServiceRouter(t *testing.T, fixture *securityFixture, script model.LLM) *mux.Router {
	t.Helper()
	service := newAIIntegrationService(t, fixture, script, "Use tools and request approval before deleting")
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
	return router
}
