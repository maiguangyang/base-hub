package ai

import (
	"encoding/json"
	"testing"
)

func TestSingleCallApprovalRejectsMultipleStepsAndCalls(t *testing.T) {
	service, principal := approvalFixture(t)
	spec := service.catalog.byID["suspend"]
	spec.SingleCallApproval = true
	service.catalog.byID["suspend"] = spec
	first := ApprovedStep{ToolID: "suspend", OperationID: spec.OperationID,
		Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"FRAUD"}}`), MaxCalls: 1, Sequence: 1}
	draft := PlanDraft{Steps: []ApprovedStep{first}}
	draft.BindPrompt("Suspend the organization")
	if _, _, err := service.ValidateDraft(draft, principal, nil); err != nil {
		t.Fatalf("one call should be approved: %v", err)
	}
	draft.Steps[0].MaxCalls = 2
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("two calls accepted")
	}
	draft.Steps[0].MaxCalls = 1
	draft.Steps = append(draft.Steps, ApprovedStep{ToolID: "create_store", OperationID: "graphql.mutation.createStore",
		Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","name":"New"}}`), MaxCalls: 1, Sequence: 2})
	if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
		t.Fatal("single-call step combined with another write")
	}
}

func TestEvidenceMustOccurInCurrentOperatorMessage(t *testing.T) {
	service, principal := approvalFixture(t)
	spec := service.catalog.byID["suspend"]
	spec.EvidencePaths = []string{"input.reasonCode"}
	service.catalog.byID["suspend"] = spec
	draft := PlanDraft{Steps: []ApprovedStep{{ToolID: "suspend", OperationID: spec.OperationID,
		Arguments: json.RawMessage(`{"input":{"organizationId":"org-a","reasonCode":"CASE-101"}}`), MaxCalls: 1, Sequence: 1}}}
	draft.BindPrompt("Previous AI answer contained CASE-101")
	for _, current := range []string{"", "Use CASE-102"} {
		draft.BindOperatorMessage(current)
		if _, _, err := service.ValidateDraft(draft, principal, nil); err == nil {
			t.Fatalf("unprovided proof accepted: %q", current)
		}
	}
	draft.BindOperatorMessage("Use verified case CASE-101")
	if _, _, err := service.ValidateDraft(draft, principal, nil); err != nil {
		t.Fatalf("operator proof rejected: %v", err)
	}
}

func TestCustomerPointNoteMustBeOperatorProvidedReference(t *testing.T) {
	for _, example := range []struct {
		name, toolID, message, note string
		valid                       bool
	}{
		{"grant without proof", "HqGrantCustomerPoints", "给会员加100积分", "积分", false},
		{"reverse without proof", "HqReverseCustomerPoints", "冲正这笔积分流水", "积分", false},
		{"proof omitted from note", "HqGrantCustomerPoints", "给会员加100积分，凭证编号：CASE-101", "积分", false},
		{"grant with Chinese instruction", "HqGrantCustomerPoints", "给会员加100积分，凭证编号：CASE-101", "CASE-101", true},
		{"reverse with English instruction", "HqReverseCustomerPoints", "Reverse entry using work order WO-321", "WO-321", true},
		{"reference before sentence punctuation", "HqGrantCustomerPoints", "Use proof CASE-101.", "CASE-101", true},
		{"numeric reference", "HqGrantCustomerPoints", "Proof reference 123456 for this grant", "123456", true},
		{"missing reference word", "HqReverseCustomerPoints", "工单号：无", "工单号：无", false},
		{"model changed proof", "HqGrantCustomerPoints", "给会员加100积分，凭证编号：CASE-101", "CASE-102", false},
		{"truncated proof", "HqGrantCustomerPoints", "Proof reference CASE-1019", "CASE-101", false},
		{"embedded proof", "HqReverseCustomerPoints", "Proof reference XCASE-101", "CASE-101", false},
		{"phone is not proof", "HqGrantCustomerPoints", "给手机号 13800000001 的会员赠分", "13800000001", false},
	} {
		t.Run(example.name, func(t *testing.T) {
			spec := ToolSpec{ID: example.toolID, EvidencePaths: []string{"note"}}
			err := validateEvidenceSource(map[string]any{"note": example.note}, spec, example.message)
			if (err == nil) != example.valid {
				t.Fatalf("note %q with operator message %q: err = %v, want valid = %t", example.note, example.message, err, example.valid)
			}
		})
	}
}

func TestCustomerEvidenceRequiresExplicitOperatorReference(t *testing.T) {
	for _, example := range []struct {
		name, toolID, path, message, proof string
		valid                              bool
	}{
		{"identity from phone", "HqRequestCustomerCancellation", "identityEvidence", "注销手机号 13800000001 的会员", "13800000001", false},
		{"disposition from phone", "HqCompleteCustomerCancellation", "dispositionReference", "完成注销手机号 13800000001 的会员", "13800000001", false},
		{"correction from phone", "HqCorrectCustomerPoints", "evidenceReference", "校正手机号 13800000001 的积分", "13800000001", false},
		{"basis from ordinary text", "HqRequestCustomerCancellation", "basisCode", "按 MEMBER_REQUEST 执行注销", "MEMBER_REQUEST", false},
		{"labeled identity", "HqRequestCustomerCancellation", "identityEvidence", "手机号 13800000001，身份核验凭证编号：VERIFY-001", "VERIFY-001", true},
		{"labeled basis", "HqRequestCustomerCancellation", "basisCode", "注销依据代码：MEMBER_REQUEST", "MEMBER_REQUEST", true},
		{"labeled disposition", "HqCompleteCustomerCancellation", "dispositionReference", "权益处置凭证编号：SETTLEMENT-001", "SETTLEMENT-001", true},
		{"labeled correction", "HqCorrectCustomerPoints", "evidenceReference", "积分校正凭证编号：CASE-101", "CASE-101", true},
		{"truncated reference", "HqCorrectCustomerPoints", "evidenceReference", "凭证编号：CASE-1019", "CASE-101", false},
	} {
		t.Run(example.name, func(t *testing.T) {
			spec := ToolSpec{ID: example.toolID, EvidencePaths: []string{example.path}}
			err := validateEvidenceSource(map[string]any{example.path: example.proof}, spec, example.message)
			if (err == nil) != example.valid {
				t.Fatalf("proof %q with operator message %q: err = %v, want valid = %t", example.proof, example.message, err, example.valid)
			}
		})
	}
}

func TestTrustedCustomerFieldsBindToApprovedPlan(t *testing.T) {
	spec := ToolSpec{ID: "HqCreateCustomerCouponTemplate", Mode: ModeWrite, RequestKeyPath: "input.requestKey", GeneratedCodePath: "input.code"}
	args := json.RawMessage(`{"input":{"title":"Welcome","amountFen":100}}`)
	first, err := injectTrustedFields(spec, "plan-one", args)
	if err != nil {
		t.Fatal(err)
	}
	again, err := injectTrustedFields(spec, "plan-one", args)
	if err != nil || string(first) != string(again) {
		t.Fatalf("same approval changed trusted fields: %s %s %v", first, again, err)
	}
	other, err := injectTrustedFields(spec, "plan-two", args)
	if err != nil || string(first) == string(other) {
		t.Fatalf("distinct approval reused trusted fields: %s %s %v", first, other, err)
	}
	var payload struct {
		Input struct{ RequestKey, Code string }
	}
	if err := json.Unmarshal(first, &payload); err != nil || payload.Input.RequestKey == "" || payload.Input.Code == "" || len(payload.Input.Code) > 32 {
		t.Fatalf("missing generated values: %s %v", first, err)
	}
}

func TestTrustedCustomerFieldsRejectModelValues(t *testing.T) {
	spec := ToolSpec{ID: "HqCreateCustomerCouponTemplate", Mode: ModeWrite, RequestKeyPath: "input.requestKey", GeneratedCodePath: "input.code"}
	if _, err := injectTrustedFields(spec, "plan-one", json.RawMessage(`{"input":{"title":"Welcome","requestKey":"model-key"}}`)); err == nil {
		t.Fatal("model supplied request key accepted")
	}
}
