package tools

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/src/services/ai"
	"google.golang.org/adk/v2/agent"
)

func TestCustomerToolsCoverExactlyReviewedOperations(t *testing.T) {
	want := map[string]struct {
		operation, permission string
		mode                  ai.ToolMode
	}{
		"HqCustomerMembers":                  {"graphql.query.hqCustomerMembers", "hqCustomer:read", ai.ModeReadOnly},
		"HqCustomerMember":                   {"graphql.query.hqCustomerMember", "hqCustomer:read", ai.ModeReadOnly},
		"HqCustomerSensitivePhone":           {"graphql.query.hqCustomerSensitivePhone", "customer:read_sensitive", ai.ModeReadOnly},
		"HqCustomerBenefitPolicy":            {"graphql.query.hqCustomerBenefitPolicy", "hqCustomerPolicy:read", ai.ModeReadOnly},
		"HqCustomerPointEntries":             {"graphql.query.hqCustomerPointEntries", "hqCustomerPoints:read", ai.ModeReadOnly},
		"HqCustomerCouponTemplates":          {"graphql.query.hqCustomerCouponTemplates", "hqCustomerCoupon:read", ai.ModeReadOnly},
		"HqCustomerCouponGrants":             {"graphql.query.hqCustomerCouponGrants", "hqCustomerCoupon:read", ai.ModeReadOnly},
		"HqCreateCustomerMember":             {"graphql.mutation.hqCreateCustomerMember", "hqCustomer:create", ai.ModeWrite},
		"HqSetCustomerMemberStatus":          {"graphql.mutation.hqSetCustomerMemberStatus", "hqCustomer:update", ai.ModeWrite},
		"HqRequestCustomerCancellation":      {"graphql.mutation.hqRequestCustomerCancellation", "hqCustomer:cancel", ai.ModeWrite},
		"HqCompleteCustomerCancellation":     {"graphql.mutation.hqCompleteCustomerCancellation", "hqCustomer:cancel", ai.ModeWrite},
		"HqSaveCustomerBenefitPolicy":        {"graphql.mutation.hqSaveCustomerBenefitPolicy", "hqCustomerPolicy:manage", ai.ModeWrite},
		"HqGrantCustomerPoints":              {"graphql.mutation.hqGrantCustomerPoints", "hqCustomerPoints:grant", ai.ModeWrite},
		"HqReverseCustomerPoints":            {"graphql.mutation.hqReverseCustomerPoints", "hqCustomerPoints:reverse", ai.ModeWrite},
		"HqCorrectCustomerPoints":            {"graphql.mutation.hqCorrectCustomerPoints", "hqCustomerPoints:correct", ai.ModeWrite},
		"HqCreateCustomerCouponTemplate":     {"graphql.mutation.hqCreateCustomerCouponTemplate", "hqCustomerCoupon:manage", ai.ModeWrite},
		"HqSetCustomerCouponTemplateEnabled": {"graphql.mutation.hqSetCustomerCouponTemplateEnabled", "hqCustomerCoupon:manage", ai.ModeWrite},
		"HqGrantCustomerCoupon":              {"graphql.mutation.hqGrantCustomerCoupon", "hqCustomerCoupon:grant", ai.ModeWrite},
		"HqRevokeCustomerCoupon":             {"graphql.mutation.hqRevokeCustomerCoupon", "hqCustomerCoupon:revoke", ai.ModeWrite},
	}
	for _, spec := range Specs() {
		if !strings.Contains(spec.OperationID, "Customer") {
			continue
		}
		expected, ok := want[spec.ID]
		if !ok {
			t.Errorf("unreviewed customer tool: %s", spec.ID)
			continue
		}
		if !customerSpecMatches(spec, expected.operation, expected.permission, expected.mode) {
			t.Errorf("invalid customer tool: %+v", spec)
		}
		delete(want, spec.ID)
	}
	if len(want) != 0 {
		t.Fatalf("missing customer tools: %v", want)
	}
}

func customerSpecMatches(spec ai.ToolSpec, operation, permission string, mode ai.ToolMode) bool {
	return spec.OperationID == operation && spec.Permission == permission && spec.Mode == mode &&
		len(spec.Workspaces) == 1 && spec.Workspaces[0] == auth.WorkspaceTypeHeadquarters && len(spec.OutputFields) > 0
}

func TestCustomerSensitivePhoneProjectsOnlyAuthorizedScalar(t *testing.T) {
	spec := customerToolSpec(t, "HqCustomerSensitivePhone")
	item, err := buildReviewedTool(spec, ai.FixedToolRuntime{
		Call: func(_ agent.Context, _ ai.ToolSpec, input json.RawMessage) (ai.FixedResponse, error) {
			if string(input) != `{"id":"member-1"}` {
				t.Fatalf("unexpected member lookup: %s", input)
			}
			return ai.FixedResponse{Status: http.StatusOK, Body: []byte(`{"data":{"hqCustomerSensitivePhone":"13800138000"}}`)}, nil
		},
		DeliverSecret: func(agent.Context, ai.SecretPayload) error {
			t.Fatal("phone was emitted as a secret event")
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	output, err := item.(interface {
		Run(agent.Context, any) (map[string]any, error)
	}).Run(toolTestContext{}, map[string]any{"id": "member-1"})
	if err != nil || !reflect.DeepEqual(output, map[string]any{"hqCustomerSensitivePhone": "13800138000"}) {
		t.Fatalf("phone projection = %v, %v", output, err)
	}
}

func TestCustomerPhoneIsNotReturnedInProtectedCallErrors(t *testing.T) {
	spec := customerToolSpec(t, "HqCustomerSensitivePhone")
	for _, response := range []ai.FixedResponse{
		{Status: http.StatusForbidden, Body: []byte(`{"error":"13800138000"}`)},
		{Status: http.StatusOK, Body: []byte(`{"data":{"hqCustomerSensitivePhone":"13800138000"},"errors":[{"message":"13800138000"}]}`)},
	} {
		item, err := buildReviewedTool(spec, ai.FixedToolRuntime{Call: func(agent.Context, ai.ToolSpec, json.RawMessage) (ai.FixedResponse, error) {
			return response, nil
		}})
		if err != nil {
			t.Fatal(err)
		}
		_, err = item.(interface {
			Run(agent.Context, any) (map[string]any, error)
		}).Run(toolTestContext{}, map[string]any{"id": "member-1"})
		if err == nil || strings.Contains(err.Error(), "13800138000") {
			t.Fatalf("protected error leaked phone: %v", err)
		}
	}
}

func customerToolSpec(t *testing.T, id string) ai.ToolSpec {
	t.Helper()
	for _, spec := range Specs() {
		if spec.ID == id {
			return spec
		}
	}
	t.Fatalf("missing customer tool %s", id)
	return ai.ToolSpec{}
}

func TestCustomerWriteSchemasHideTrustedFields(t *testing.T) {
	for _, example := range []struct {
		id, path string
	}{
		{"HqCreateCustomerMember", "input.requestKey"},
		{"HqCreateCustomerCouponTemplate", "input.requestKey"},
		{"HqCreateCustomerCouponTemplate", "input.code"},
		{"HqGrantCustomerPoints", "requestKey"},
		{"HqReverseCustomerPoints", "requestKey"},
		{"HqCorrectCustomerPoints", "requestKey"},
		{"HqGrantCustomerCoupon", "requestKey"},
	} {
		var found bool
		for _, spec := range Specs() {
			if spec.ID != example.id {
				continue
			}
			found = true
			schema, err := reviewedInputSchema(spec)
			if err != nil {
				t.Fatal(err)
			}
			parts := strings.Split(example.path, ".")
			for _, part := range parts[:len(parts)-1] {
				schema = schema.Properties[part]
			}
			if _, exposed := schema.Properties[parts[len(parts)-1]]; exposed {
				t.Errorf("%s exposes trusted field %s", example.id, example.path)
			}
		}
		if !found {
			t.Fatalf("missing tool %s", example.id)
		}
	}
}

func TestEveryCustomerWriteRequiresSingleCallApproval(t *testing.T) {
	specs := customerWriteSpecs()
	if len(specs) != 12 {
		t.Fatalf("customer write tools = %d", len(specs))
	}
	for _, spec := range specs {
		if spec.Mode != ai.ModeWrite || !spec.SingleCallApproval {
			t.Errorf("customer write lacks one-call confirmation: %s", spec.ID)
		}
	}
}

func TestCustomerMemberCreateToolOnlyAsksForPhone(t *testing.T) {
	spec := customerToolSpec(t, "HqCreateCustomerMember")
	if len(spec.EvidencePaths) != 0 {
		t.Fatalf("member creation still asks for evidence: %v", spec.EvidencePaths)
	}
	schema, err := reviewedInputSchema(spec)
	if err != nil {
		t.Fatal(err)
	}
	fields := schema.Properties["input"].Properties
	if len(fields) != 1 || fields["phone"] == nil {
		t.Fatalf("member creation AI fields = %v", fields)
	}
	if !strings.Contains(fields["phone"].Description, "11 位") {
		t.Fatalf("member creation phone guidance = %q", fields["phone"].Description)
	}
}

func TestCustomerMemberCreateApprovalRejectsInvalidPhones(t *testing.T) {
	specs, err := PreparedSpecs()
	if err != nil {
		t.Fatal(err)
	}
	for _, spec := range specs {
		if spec.ID != "HqCreateCustomerMember" {
			continue
		}
		for _, phone := range []string{"1380000000", "138000000000", "12800000000", "1380000000a", "+8613800000001", "138-0000-0001"} {
			if err := spec.InputSchema.Validate(map[string]any{"input": map[string]any{"phone": phone}}); err == nil {
				t.Errorf("approval accepted invalid phone %q", phone)
			}
		}
		if err := spec.InputSchema.Validate(map[string]any{"input": map[string]any{"phone": "13800000001"}}); err != nil {
			t.Fatalf("approval rejected valid phone: %v", err)
		}
		return
	}
	t.Fatal("missing member creation tool")
}

func TestCustomerPointNotesRequireOperatorSource(t *testing.T) {
	for _, id := range []string{"HqGrantCustomerPoints", "HqReverseCustomerPoints"} {
		spec := customerToolSpec(t, id)
		if !reflect.DeepEqual(spec.EvidencePaths, []string{"note"}) {
			t.Errorf("%s note is not bound to operator message: %v", id, spec.EvidencePaths)
		}
	}
}

func TestCustomerPointNoteSchemaExplainsReferenceOnly(t *testing.T) {
	for _, id := range []string{"HqGrantCustomerPoints", "HqReverseCustomerPoints"} {
		schema, err := reviewedInputSchema(customerToolSpec(t, id))
		if err != nil {
			t.Fatal(err)
		}
		description := schema.Properties["note"].Description
		if !strings.Contains(description, "仅填写凭证编号") {
			t.Errorf("%s note description does not explain reference-only input: %q", id, description)
		}
	}
}

func TestCustomerStatusToolOnlyAcceptsActiveAndSuspended(t *testing.T) {
	spec := customerToolSpec(t, "HqSetCustomerMemberStatus")
	schema, err := reviewedInputSchema(spec)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := schema.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"ACTIVE", "SUSPENDED"} {
		if err := resolved.Validate(map[string]any{"id": "member-1", "status": status}); err != nil {
			t.Errorf("valid status %s rejected: %v", status, err)
		}
	}
	for _, status := range []string{"CANCEL_PENDING", "CANCELLED"} {
		if err := resolved.Validate(map[string]any{"id": "member-1", "status": status}); err == nil {
			t.Errorf("invalid status %s accepted", status)
		}
	}
}
