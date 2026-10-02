package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"base-engine/gen"
	"base-engine/src/services/ai"
	aitools "base-engine/src/services/ai/tools"
	"google.golang.org/adk/v2/agent"
)

func TestCouponAIToolsScheduleMemberCatchupOnlyOnCreateAndReactivation(t *testing.T) {
	fixture := newCustomerFixture(t)
	created, err := runCouponAITool(t, fixture, "HqCreateCustomerMember", map[string]any{
		"input": map[string]any{"phone": "13800000009"},
	}, "member-create")
	if err != nil {
		t.Fatal(err)
	}
	memberID := couponAIResultID(t, created, "hqCreateCustomerMember")

	var jobs []gen.CustomerCouponDistributionJob
	if err := fixture.db.Where("member_id = ?", memberID).Find(&jobs).Error; err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].Kind != gen.CustomerCouponDistributionJobKindMemberCatchup {
		t.Fatalf("member catch-up jobs = %#v", jobs)
	}
	if err := fixture.db.Model(&gen.CustomerCouponDistributionJob{}).Where("id = ?", jobs[0].ID).
		Updates(map[string]any{"status": gen.CustomerCouponDistributionJobStatusCompleted, "attempts": 3}).Error; err != nil {
		t.Fatal(err)
	}

	if _, err := runCouponAITool(t, fixture, "HqSetCustomerMemberStatus", map[string]any{"id": memberID, "status": "ACTIVE"}, "active-noop"); err != nil {
		t.Fatal(err)
	}
	assertCouponAIJobState(t, fixture, jobs[0].ID, gen.CustomerCouponDistributionJobStatusCompleted, 3)

	if _, err := runCouponAITool(t, fixture, "HqSetCustomerMemberStatus", map[string]any{"id": memberID, "status": "SUSPENDED"}, "suspend"); err != nil {
		t.Fatal(err)
	}
	if _, err := runCouponAITool(t, fixture, "HqSetCustomerMemberStatus", map[string]any{"id": memberID, "status": "ACTIVE"}, "reactivate"); err != nil {
		t.Fatal(err)
	}
	assertCouponAIJobState(t, fixture, jobs[0].ID, gen.CustomerCouponDistributionJobStatusPending, 0)
}

func TestCouponAIToolsRejectFutureManualGrantWithoutWriting(t *testing.T) {
	fixture := newCustomerFixture(t)
	memberID := createCouponAIMember(t, fixture, "13800000008", "future-grant-member")
	future := time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
	template, err := runCouponAITool(t, fixture, "HqCreateCustomerCouponTemplate", map[string]any{
		"input": map[string]any{
			"title": "Future", "amountFen": 500, "minSpendFen": 0,
			"daysAfterActivation": 30, "effectiveAt": future,
			"perMemberLimit": 1, "enabled": true,
		},
	}, "future-template")
	if err != nil {
		t.Fatal(err)
	}
	templateID := couponAIResultID(t, template, "hqCreateCustomerCouponTemplate")

	if output, err := runCouponAITool(t, fixture, "HqGrantCustomerCoupon", map[string]any{
		"templateId": templateID, "memberId": memberID,
	}, "future-grant"); err == nil {
		t.Fatalf("future manual grant succeeded: %#v", output)
	}
	var count int64
	if err := fixture.db.Model(&gen.CustomerCouponGrant{}).
		Where("template_id = ? AND member_id = ?", templateID, memberID).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("rejected manual grant wrote %d rows", count)
	}
}

func TestCouponAIToolsRejectCrossStoreManualGrantWithoutWriting(t *testing.T) {
	fixture := newCustomerFixture(t)
	memberID := createCouponAIMember(t, fixture, "13800000007", "cross-store-grant-member")
	permission := gen.Permission{ID: "coupon-franchise-grant", Action: "franchiseCoupon:grant", Name: "franchiseCoupon:grant", Module: "customer", Scope: gen.PermissionScopeTenant}
	secondStore := gen.Store{ID: "store-a-2", Code: "A2", Name: "A Store 2", Lifecycle: gen.StoreLifecycleActive, OrganizationID: "org-a"}
	templateStoreID := "store-a"
	storeTemplate := gen.CustomerCouponTemplate{
		ID: "store-template-a", Code: "STOREA", RequestKey: "store-template-a", Title: "Store A",
		AmountFen: 500, DaysAfterActivation: 30, EffectiveAt: time.Now().Add(-time.Hour).UnixMilli(),
		PerMemberLimit: 1, Enabled: true, IssuerScope: gen.CouponIssuerScopeStore,
		OrganizationID: "org-hq", ApplicableStoreID: &templateStoreID,
	}
	fixture.createAll([]gen.Permission{permission}, []permissionRole{{PermissionID: permission.ID, OperatorRoleID: "role-owner-a"}}, []gen.Store{secondStore}, []gen.CustomerCouponTemplate{storeTemplate})
	if output, err := runCouponAIToolAs(t, fixture, "session-a-1", "FranchiseGrantCoupon", map[string]any{
		"storeId": secondStore.ID, "templateId": storeTemplate.ID, "memberId": memberID,
	}, "wrong-store-grant"); err == nil {
		t.Fatalf("cross-store manual grant succeeded: %#v", output)
	}
	var count int64
	if err := fixture.db.Model(&gen.CustomerCouponGrant{}).
		Where("template_id = ? AND member_id = ?", storeTemplate.ID, memberID).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("rejected cross-store grant wrote %d rows", count)
	}
}

func createCouponAIMember(t *testing.T, fixture *securityFixture, phone, key string) string {
	t.Helper()
	member, err := runCouponAITool(t, fixture, "HqCreateCustomerMember", map[string]any{
		"input": map[string]any{"phone": phone},
	}, key)
	if err != nil {
		t.Fatal(err)
	}
	return couponAIResultID(t, member, "hqCreateCustomerMember")
}

func runCouponAITool(t *testing.T, fixture *securityFixture, toolID string, args map[string]any, trustedSuffix string) (map[string]any, error) {
	return runCouponAIToolAs(t, fixture, "session-hq", toolID, args, trustedSuffix)
}

func runCouponAIToolAs(t *testing.T, fixture *securityFixture, sessionID, toolID string, args map[string]any, trustedSuffix string) (map[string]any, error) {
	t.Helper()
	runtime := ai.FixedToolRuntime{Call: func(_ agent.Context, spec ai.ToolSpec, variables json.RawMessage) (ai.FixedResponse, error) {
		var values map[string]any
		if err := json.Unmarshal(variables, &values); err != nil {
			return ai.FixedResponse{}, err
		}
		requestKey := fmt.Sprintf("ai-coupon-%s-%s", trustedSuffix, spec.ID)
		if input, ok := values["input"].(map[string]any); ok {
			if spec.RequestKeyPath == "input.requestKey" {
				input["requestKey"] = requestKey
			}
			if spec.GeneratedCodePath == "input.code" {
				input["code"] = "AI" + fmt.Sprintf("%x", requestKey)[:24]
			}
		}
		if spec.RequestKeyPath == "requestKey" {
			values["requestKey"] = requestKey
		}
		encoded, err := json.Marshal(values)
		if err != nil {
			return ai.FixedResponse{}, err
		}
		return ai.ProtectedCall(t.Context(), fixture.handler, ai.FixedRequest{Method: http.MethodPost, Path: "/graphql", Document: spec.Document,
			Variables: encoded, Cookie: fixture.cookies[sessionID].String(), Origin: "https://admin.example.com"})
	}, DeliverSecret: func(agent.Context, ai.SecretPayload) error { return nil }}
	items, err := aitools.Build(runtime)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.Name() == toolID {
			return item.(interface {
				Run(agent.Context, any) (map[string]any, error)
			}).Run(protectedGapToolContext{}, args)
		}
	}
	t.Fatalf("missing tool %s", toolID)
	return nil, nil
}

func couponAIResultID(t *testing.T, result map[string]any, root string) string {
	t.Helper()
	value, ok := result[root].(map[string]any)
	if !ok {
		t.Fatalf("%s result = %#v", root, result)
	}
	id, _ := value["id"].(string)
	if id == "" {
		t.Fatalf("%s missing id: %#v", root, value)
	}
	return id
}

func assertCouponAIJobState(t *testing.T, fixture *securityFixture, id string, status gen.CustomerCouponDistributionJobStatus, attempts int64) {
	t.Helper()
	var job gen.CustomerCouponDistributionJob
	if err := fixture.db.First(&job, "id = ?", id).Error; err != nil {
		t.Fatal(err)
	}
	if job.Status != status || job.Attempts != attempts {
		t.Fatalf("job state = %s/%d, want %s/%d", job.Status, job.Attempts, status, attempts)
	}
}
