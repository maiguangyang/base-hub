package customer

import (
	"context"
	"testing"

	"base-engine/gen"
)

func TestCouponFiltersApplyBeforePagination(t *testing.T) {
	service, _, principal, _ := couponFixture(t)
	ctx := context.Background()
	first, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil {
		t.Fatal(err)
	}
	other := validCouponInput()
	other.Code, other.Title, other.RequestKey, other.Enabled = "SECOND", "Second", "template-two", false
	if _, err := service.CreateCouponTemplate(ctx, principal, other); err != nil {
		t.Fatal(err)
	}
	q, enabled := "WELCOME", true
	page, err := service.ListCouponTemplatesFiltered(ctx, principal, CouponTemplateFilter{Q: &q, Enabled: &enabled}, 1, 1)
	if err != nil || page.Total != 1 || len(page.Data) != 1 || page.Data[0].ID != first.ID {
		t.Fatalf("filtered templates = %+v, %v", page, err)
	}
}

func TestCouponGrantStatusFiltersApplyBeforePagination(t *testing.T) {
	service, _, principal, member := couponFixture(t)
	ctx := context.Background()
	first, err := service.CreateCouponTemplate(ctx, principal, validCouponInput())
	if err != nil { t.Fatal(err) }
	if _, err := service.GrantCoupon(ctx, principal, first.ID, member.ID, "filter-grant"); err != nil {
		t.Fatal(err)
	}
	status := gen.CustomerCouponGrantStatusPendingActivation
	grants, err := service.ListCouponGrantsFiltered(ctx, principal, member.ID, CouponGrantFilter{TemplateID: &first.ID, Status: &status}, 1, 1)
	if err != nil || grants.Total != 1 || len(grants.Data) != 1 {
		t.Fatalf("filtered grants = %+v, %v", grants, err)
	}
	status = gen.CustomerCouponGrantStatusRevoked
	grants, err = service.ListCouponGrantsFiltered(ctx, principal, member.ID, CouponGrantFilter{Status: &status}, 1, 1)
	if err != nil || grants.Total != 0 {
		t.Fatalf("revoked grants = %+v, %v", grants, err)
	}
}
