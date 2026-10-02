package integration_test

import (
	"encoding/json"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func newCustomerFixture(t *testing.T) *securityFixture {
	t.Helper()
	fixture := newSecurityFixture(t)
	for _, model := range []any{&gen.CustomerMember{}, &gen.CustomerBenefitPolicy{}, &gen.CustomerDailyPointGrantBudget{}, &gen.CustomerPointEntry{}, &gen.CustomerCouponTemplate{}, &gen.CustomerCouponGrant{}, &gen.CustomerCouponDistributionJob{}} {
		if err := fixture.db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			fixture.db.Exec("DROP INDEX IF EXISTS `" + index + "`")
		}
	}
	if err := fixture.db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS uk_customer_coupon_distribution_jobs_request_key ON customer_coupon_distribution_jobs(request_key)").Error; err != nil {
		t.Fatal(err)
	}
	for _, action := range []string{"hqCustomer:read", "hqCustomer:create", "hqCustomer:update", "hqCustomer:cancel", "customer:read_sensitive", "hqCustomerPolicy:read", "hqCustomerPolicy:manage", "hqCustomerPoints:read", "hqCustomerPoints:grant", "hqCustomerPoints:reverse", "hqCustomerPoints:correct", "hqCustomerCoupon:read", "hqCustomerCoupon:manage", "hqCustomerCoupon:grant", "hqCustomerCoupon:revoke"} {
		id := "customer-" + strings.ReplaceAll(action, ":", "-")
		fixture.createAll([]gen.Permission{{ID: id, Action: action, Name: action, Module: "customer", Scope: gen.PermissionScopeSystem}},
			[]permissionRole{{PermissionID: id, OperatorRoleID: "role-hq"}})
	}
	return fixture
}

func TestCustomerCouponGraphQL(t *testing.T) {
	fixture := newCustomerFixture(t)
	create := `mutation { hqCreateCustomerCouponTemplate(input:{code:"WELCOME", requestKey:"template-one", title:"Welcome", amountFen:500, minSpendFen:1000, daysAfterActivation:30, perMemberLimit:1, totalIssueLimit:1, enabled:true}) { id code enabled } }`
	assertCode(t, fixture.execute("session-a-1", create), auth.CodePermissionDenied)
	templateResponse := fixture.execute("session-hq", create)
	assertCustomerSuccess(t, templateResponse)
	var template struct{ ID string }
	if err := json.Unmarshal(templateResponse.Data["hqCreateCustomerCouponTemplate"], &template); err != nil {
		t.Fatal(err)
	}
	memberResponse := fixture.execute("session-hq", `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"member-one"}) { id } }`)
	assertCustomerSuccess(t, memberResponse)
	var member struct{ ID string }
	if err := json.Unmarshal(memberResponse.Data["hqCreateCustomerMember"], &member); err != nil {
		t.Fatal(err)
	}
	grant := `mutation { hqGrantCustomerCoupon(templateId:"` + template.ID + `", memberId:"` + member.ID + `", requestKey:"grant-1") { id status activatedAt expiresAt } }`
	granted := fixture.execute("session-hq", grant)
	assertCustomerSuccess(t, granted)
	if !strings.Contains(granted.Body, "PENDING_ACTIVATION") || !strings.Contains(granted.Body, `"expiresAt":null`) {
		t.Fatalf("pending grant = %s", granted.Body)
	}
	assertCode(t, fixture.execute("session-hq", `mutation { updateCustomerCouponTemplate(id:"`+template.ID+`", input:{amountFen:1}) { id } }`), auth.CodePermissionDenied)
}

func TestCustomerPointGraphQL(t *testing.T) {
	fixture := newCustomerFixture(t)
	policy := gen.CustomerBenefitPolicy{
		ID: "policy-hq", OrganizationID: "org-hq", Version: 1,
		ManualGrantMaxSingle: 100, ManualGrantMaxDaily: 100,
	}
	fixture.createAll([]gen.CustomerBenefitPolicy{policy})
	create := fixture.execute("session-hq", `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"member-one"}) { id } }`)
	assertCustomerSuccess(t, create)
	var member struct{ ID string }
	if err := json.Unmarshal(create.Data["hqCreateCustomerMember"], &member); err != nil {
		t.Fatal(err)
	}
	gift := `mutation { hqGrantCustomerPoints(memberId:"` + member.ID + `", points:50, reasonCode:REWARD, note:"gift", requestKey:"gift-1") { entry { id delta } currentBalance } }`
	assertCode(t, fixture.execute("session-a-1", gift), auth.CodePermissionDenied)
	result := fixture.execute("session-hq", gift)
	assertCustomerSuccess(t, result)
	if !strings.Contains(result.Body, `"currentBalance":50`) {
		t.Fatalf("gift result = %s", result.Body)
	}
	list := fixture.execute("session-hq", `query { hqCustomerPointEntries(memberId:"`+member.ID+`", page:1, perPage:20) { total data { id delta } } }`)
	assertCustomerSuccess(t, list)
	if !strings.Contains(list.Body, `"total":1`) {
		t.Fatalf("point history = %s", list.Body)
	}
	assertCode(t, fixture.execute("session-hq", `mutation { updateCustomerMember(id:"`+member.ID+`", input:{pointsBalance:500}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-hq", `query { customerPointEntries { total } }`), auth.CodePermissionDenied)
}

func TestCustomerPolicyGraphQL(t *testing.T) {
	fixture := newCustomerFixture(t)
	query := `query { hqCustomerBenefitPolicy { version redemptionEnabled redeemPoints redeemAmountFen } }`
	assertCode(t, fixture.execute("session-a-1", query), auth.CodePermissionDenied)
	draft := fixture.execute("session-hq", query)
	assertCustomerSuccess(t, draft)
	if !strings.Contains(draft.Body, `"version":0`) || !strings.Contains(draft.Body, `"redeemPoints":100`) {
		t.Fatalf("policy draft: %s", draft.Body)
	}
	save := `mutation { hqSaveCustomerBenefitPolicy(expectedVersion:0, input:{discountEnabled:true, discountBasisPoints:9500, purchaseEarnEnabled:true, earnAmountFen:100, earnPoints:1, redemptionEnabled:true, redeemPoints:100, redeemAmountFen:100, maxRedemptionBasisPoints:5000, maxRedemptionPoints:5000, manualGrantMaxSingle:0, manualGrantMaxDaily:0}) { version discountBasisPoints manualGrantMaxSingle manualGrantMaxDaily } }`
	saved := fixture.execute("session-hq", save)
	assertCustomerSuccess(t, saved)
	if !strings.Contains(saved.Body, `"version":1`) || !strings.Contains(saved.Body, `"manualGrantMaxSingle":0`) || !strings.Contains(saved.Body, `"manualGrantMaxDaily":0`) {
		t.Fatalf("saved policy: %s", saved.Body)
	}
	assertCode(t, fixture.execute("session-hq", save), auth.CodeConflict)
}

func assertCustomerSuccess(t *testing.T, response graphQLResponse) {
	t.Helper()
	if len(response.Errors) != 0 {
		t.Fatalf("GraphQL request failed: %s", response.Body)
	}
}

func TestCustomerGraphQLRelationIDsDenied(t *testing.T) {
	fixture := newCustomerFixture(t)
	member := gen.CustomerMember{
		ID: "customer-one", MemberNumber: "CM-one", RequestKey: "request-one",
		Status: gen.CustomerMemberStatusActive, OrganizationID: "org-hq",
		NoticeVersion: "v1", ProcessingBasisCode: "service",
		EvidenceReference: "evidence",
	}
	if err := fixture.db.Create(&member).Error; err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"customerMembersIds", "customerBenefitPoliciesIds", "customerDailyPointGrantBudgetsIds", "customerPointEntriesIds", "customerCouponTemplatesIds"} {
		t.Run(field, func(t *testing.T) {
			response := fixture.execute("session-hq", `query { organization(id: "org-hq") { `+field+` } }`)
			assertCode(t, response, auth.CodePermissionDenied)
		})
	}
}

func TestCustomerGraphQLMemberFlow(t *testing.T) {
	fixture := newCustomerFixture(t)
	create := `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"create-1"}) { id memberNumber phoneMasked status } }`
	assertCode(t, fixture.execute("session-a-1", create), auth.CodePermissionDenied)
	response := fixture.execute("session-hq", create)
	assertCustomerSuccess(t, response)
	var created struct {
		ID, MemberNumber, PhoneMasked, Status string
	}
	if err := json.Unmarshal(response.Data["hqCreateCustomerMember"], &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.MemberNumber == "" {
		t.Fatalf("missing created identity: %+v", created)
	}
	if created.PhoneMasked != "+86138****0001" || created.Status != "ACTIVE" {
		t.Fatalf("created = %+v", created)
	}
	list := fixture.execute("session-hq", `query { hqCustomerMembers(phone:"1380000", page:1, perPage:20) { total data { id phoneMasked } } }`)
	assertCustomerSuccess(t, list)
	if !strings.Contains(list.Body, created.ID) || strings.Contains(list.Body, "+8613800000001") {
		t.Fatalf("unsafe list: %s", list.Body)
	}
	assertCode(t, fixture.execute("session-a-1", `query { hqCustomerMember(id:"`+created.ID+`") { id } }`), auth.CodePermissionDenied)
	phone := fixture.execute("session-hq", `query { hqCustomerSensitivePhone(id:"`+created.ID+`") }`)
	assertCustomerSuccess(t, phone)
	if !strings.Contains(phone.Body, "+8613800000001") {
		t.Fatalf("sensitive phone = %s", phone.Body)
	}
	assertCode(t, fixture.execute("session-hq", `query { customerMember(id:"`+created.ID+`") { phone } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-hq", `mutation { updateCustomerMember(id:"`+created.ID+`", input:{pointsBalance:100}) { id } }`), auth.CodePermissionDenied)
	assertCode(t, fixture.execute("session-hq", `mutation { deleteCustomerMembers(id:["`+created.ID+`"]) }`), auth.CodePermissionDenied)
}

func TestCustomerGraphQLCancellation(t *testing.T) {
	fixture := newCustomerFixture(t)
	create := fixture.execute("session-hq", `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"create-1"}) { id } }`)
	assertCustomerSuccess(t, create)
	var created struct{ ID string }
	if err := json.Unmarshal(create.Data["hqCreateCustomerMember"], &created); err != nil {
		t.Fatal(err)
	}
	assertCode(t, fixture.execute("session-hq", `mutation { hqSetCustomerMemberStatus(id:"`+created.ID+`", status:CANCELLED) { status } }`), auth.CodeValidationFailed)
	request := fixture.execute("session-hq", `mutation { hqRequestCustomerCancellation(id:"`+created.ID+`", identityEvidence:"identity-1", basisCode:"REQUEST") { status } }`)
	assertCustomerSuccess(t, request)
	if !strings.Contains(request.Body, "CANCEL_PENDING") {
		t.Fatalf("cancellation request: %s", request.Body)
	}
	complete := fixture.execute("session-hq", `mutation { hqCompleteCustomerCancellation(id:"`+created.ID+`", dispositionReference:"rights-settled") { status } }`)
	assertCustomerSuccess(t, complete)
	if !strings.Contains(complete.Body, "CANCELLED") {
		t.Fatalf("cancellation completion: %s", complete.Body)
	}
	var member gen.CustomerMember
	if err := fixture.db.First(&member, "id = ?", created.ID).Error; err != nil || member.Phone != nil {
		t.Fatalf("cancelled phone = %v, %v", member.Phone, err)
	}
	rejoin := fixture.execute("session-hq", `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"create-2"}) { id } }`)
	assertCustomerSuccess(t, rejoin)
	if strings.Contains(rejoin.Body, created.ID) {
		t.Fatalf("new member after cancellation: %s", rejoin.Body)
	}
}

func TestCustomerGraphQLMutationRightsCount(t *testing.T) {
	fixture := newCustomerFixture(t)
	create := fixture.execute("session-hq", `mutation { hqCreateCustomerMember(input: {phone:"13800000001", requestKey:"create-1"}) { id } }`)
	assertCustomerSuccess(t, create)
	var member struct{ ID string }
	if err := json.Unmarshal(create.Data["hqCreateCustomerMember"], &member); err != nil {
		t.Fatal(err)
	}
	createCustomerPendingGrant(t, fixture, member.ID)
	status := fixture.execute("session-hq", `mutation { hqSetCustomerMemberStatus(id:"`+member.ID+`", status:SUSPENDED) { status pendingCouponCount } }`)
	assertCustomerSuccess(t, status)
	if !strings.Contains(status.Body, `"pendingCouponCount":1`) {
		t.Fatalf("status rights = %s", status.Body)
	}
	request := fixture.execute("session-hq", `mutation { hqRequestCustomerCancellation(id:"`+member.ID+`", identityEvidence:"verify", basisCode:"REQUEST") { status pendingCouponCount } }`)
	assertCustomerSuccess(t, request)
	if !strings.Contains(request.Body, `"pendingCouponCount":1`) {
		t.Fatalf("cancellation rights = %s", request.Body)
	}
}

func createCustomerPendingGrant(t *testing.T, fixture *securityFixture, memberID string) {
	t.Helper()
	template := gen.CustomerCouponTemplate{
		ID: "template-rights", OrganizationID: "org-hq", Code: "RIGHTS", RequestKey: "template-rights",
		Title: "Rights", AmountFen: 500, DaysAfterActivation: 30, PerMemberLimit: 1, TotalIssueLimit: 1,
	}
	grant := gen.CustomerCouponGrant{
		ID: "grant-rights", MemberID: memberID, TemplateID: template.ID,
		Status: gen.CustomerCouponGrantStatusPendingActivation, RequestKey: "grant-rights",
		AmountFen: 500, DaysAfterActivation: 30,
	}
	fixture.createAll([]gen.CustomerCouponTemplate{template}, []gen.CustomerCouponGrant{grant})
}
