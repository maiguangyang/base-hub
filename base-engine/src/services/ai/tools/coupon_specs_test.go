package tools

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
)

func TestCouponCreateSchemasExposeOptionalScheduleAndBoundedLifetime(t *testing.T) {
	for _, id := range []string{"HqCreateCustomerCouponTemplate", "FranchiseCreateCouponTemplate"} {
		spec := customerToolSpec(t, id)
		schema, err := reviewedInputSchema(spec)
		if err != nil {
			t.Fatal(err)
		}
		input := schema.Properties["input"]
		if input == nil {
			t.Fatalf("%s input schema missing", id)
		}
		assertCouponOptionalFields(t, id, input)
		days := input.Properties["daysAfterActivation"]
		if days == nil || days.Minimum == nil || *days.Minimum != 1 || days.Maximum == nil || *days.Maximum != 365 {
			t.Errorf("%s days bounds = %#v", id, days)
		}
	}
}

func assertCouponOptionalFields(t *testing.T, id string, input *jsonschema.Schema) {
	t.Helper()
	for _, optional := range []string{"effectiveAt", "distributionEndsAt", "totalIssueLimit"} {
		field := input.Properties[optional]
		if field == nil {
			t.Errorf("%s missing %s", id, optional)
			continue
		}
		if slices.Contains(input.Required, optional) {
			t.Errorf("%s unexpectedly requires %s", id, optional)
		}
		if !strings.Contains(field.Description, "省略") {
			t.Errorf("%s %s does not explain omission: %q", id, optional, field.Description)
		}
		if !slices.Contains(field.Types, "null") {
			t.Errorf("%s %s does not accept explicit null: %#v", id, optional, field)
		}
	}
}

func TestCouponCreateAndListToolsDescribeDistributionContract(t *testing.T) {
	createExpectations := map[string][]string{
		"HqCreateCustomerCouponTemplate": {"现有会员", "后注册会员", "异步"},
		"FranchiseCreateCouponTemplate":  {"现有会员", "后注册会员", "异步", "创建门店"},
	}
	for id, terms := range createExpectations {
		spec := customerToolSpec(t, id)
		for _, term := range terms {
			if !strings.Contains(spec.Description, term) {
				t.Errorf("%s description missing %q: %q", id, term, spec.Description)
			}
		}
		for _, field := range []string{"effectiveAt", "distributionEndsAt", "totalIssueLimit"} {
			if !strings.Contains(spec.Document, field) {
				t.Errorf("%s output does not select %s", id, field)
			}
		}
	}

	for _, id := range []string{"HqCustomerCouponTemplates", "FranchiseCouponTemplates"} {
		document := customerToolSpec(t, id).Document
		for _, field := range []string{"effectiveAt", "distributionEndsAt", "daysAfterActivation", "perMemberLimit", "totalIssueLimit"} {
			if !strings.Contains(document, field) {
				t.Errorf("%s list does not select %s", id, field)
			}
		}
	}
}

func TestCouponTimeAndCapacityDescriptionsAreExplicit(t *testing.T) {
	checks := map[string][]string{
		"effectiveAt":         {"省略", "立即生效"},
		"distributionEndsAt":  {"省略", "不限制", "到达该时刻不再"},
		"daysAfterActivation": {"1 至 365"},
		"totalIssueLimit":     {"省略", "不限制"},
	}
	for field, terms := range checks {
		description := inputFieldDescriptions[field]
		for _, term := range terms {
			if !strings.Contains(description, term) {
				t.Errorf("%s description missing %q: %q", field, term, description)
			}
		}
	}
}

func TestCouponMemberAndManualGrantToolsDescribeEligibilityBoundaries(t *testing.T) {
	checks := map[string][]string{
		"HqCreateCustomerMember":    {"符合条件", "异步", "总部券", "门店券"},
		"HqSetCustomerMemberStatus": {"SUSPENDED", "ACTIVE", "重新", "异步"},
		"HqGrantCustomerCoupon":     {"生效", "截止", "总发放量"},
		"FranchiseGrantCoupon":      {"生效", "截止", "总发放量", "创建门店"},
	}
	for id, terms := range checks {
		description := customerToolSpec(t, id).Description
		for _, term := range terms {
			if !strings.Contains(description, term) {
				t.Errorf("%s description missing %q: %q", id, term, description)
			}
		}
	}
}
