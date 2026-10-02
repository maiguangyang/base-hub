package customer

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
)

func TestCustomerMemberSearchMasksPhone(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	input := customerCreateInput("first")
	if _, err := service.CreateMember(context.Background(), principal, input); err != nil {
		t.Fatal(err)
	}
	input.Phone, input.RequestKey = "13900000002", "second"
	if _, err := service.CreateMember(context.Background(), principal, input); err != nil {
		t.Fatal(err)
	}
	page, err := service.SearchMembers(context.Background(), principal, MemberSearch{Phone: "1380000", Page: 1, PerPage: 20})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 1 || len(page.Data) != 1 ||
		page.Data[0].PhoneMasked != "+86138****0001" {
		t.Fatalf("masked search = %+v", page)
	}
	status := gen.CustomerMemberStatusActive
	page, err = service.SearchMembers(context.Background(), principal, MemberSearch{Status: &status, Page: 1, PerPage: 1})
	if err != nil || page.Total != 2 || len(page.Data) != 1 {
		t.Fatalf("status page = %+v, %v", page, err)
	}
}

func TestCustomerMemberSearchValidation(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	for _, search := range []MemberSearch{{Phone: "138000", Page: 1, PerPage: 20}, {Page: 1, PerPage: 51}} {
		if _, err := service.SearchMembers(context.Background(), principal, search); auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Fatalf("invalid search %+v: %v", search, err)
		}
	}
}

func TestCustomerMemberSearchIncludesPendingCouponCount(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	member, err := service.CreateMember(context.Background(), principal, customerCreateInput("first"))
	if err != nil {
		t.Fatal(err)
	}
	template := gen.CustomerCouponTemplate{
		ID: "template-1", OrganizationID: "hq", Code: "WELCOME", RequestKey: "template-1",
		Title: "Welcome", AmountFen: 500, DaysAfterActivation: 30, EffectiveAt: 1_800_000_000_000,
		PerMemberLimit: 2, TotalIssueLimit: 2,
	}
	if err := db.Create(&template).Error; err != nil {
		t.Fatal(err)
	}
	for _, item := range []struct {
		id     string
		status gen.CustomerCouponGrantStatus
	}{
		{"pending", gen.CustomerCouponGrantStatusPendingActivation},
		{"revoked", gen.CustomerCouponGrantStatusRevoked},
	} {
		grant := gen.CustomerCouponGrant{
			ID: item.id, MemberID: member.ID, TemplateID: template.ID, RequestKey: item.id,
			Status: item.status, AmountFen: 500, DaysAfterActivation: 30, IssuedAt: 1_800_000_000_000,
		}
		if err := db.Create(&grant).Error; err != nil {
			t.Fatal(err)
		}
	}
	page, err := service.SearchMembers(context.Background(), principal, MemberSearch{Page: 1, PerPage: 20})
	if err != nil || len(page.Data) != 1 || page.Data[0].PendingCouponCount != 1 {
		t.Fatalf("pending coupon count = %+v, %v", page, err)
	}
}

func TestCustomerSensitivePhoneNeedsSeparateAction(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	member, err := service.CreateMember(context.Background(), principal, customerCreateInput("first"))
	if err != nil {
		t.Fatal(err)
	}
	phone, err := service.SensitivePhone(context.Background(), principal, member.ID)
	if err != nil || phone != "+8613800000001" {
		t.Fatalf("sensitive phone = %q, %v", phone, err)
	}
	delete(principal.Permissions, "customer:read_sensitive")
	if _, err := service.SensitivePhone(context.Background(), principal, member.ID); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("missing sensitive action = %v", err)
	}
}

func TestCustomerSensitivePhonePreservesDatabaseErrors(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	connection, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = service.SensitivePhone(context.Background(), principal, "member-1")
	if err == nil || auth.ErrorCode(err) == auth.CodePermissionDenied {
		t.Fatalf("database error hidden as permission failure: %v", err)
	}
}
