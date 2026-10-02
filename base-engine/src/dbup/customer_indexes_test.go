package dbup

import (
	"testing"

	"base-engine/gen"
	"gorm.io/gorm"
)

func TestCustomerIndexesProtectHQIdentity(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGeneratedForSQLite(db,
		&gen.Organization{}, &gen.CustomerMember{}, &gen.CustomerBenefitPolicy{},
		&gen.CustomerDailyPointGrantBudget{}, &gen.CustomerPointEntry{},
		&gen.CustomerCouponTemplate{}, &gen.CustomerCouponGrant{}, &gen.CustomerCouponDistributionJob{},
	); err != nil {
		t.Fatal(err)
	}
	if err := EnsureCustomerIndexes(db); err != nil {
		t.Fatal(err)
	}
	for _, item := range []struct{ table, index string }{
		{"customer_members", "uidx_customer_member_organization_phone"},
		{"customer_members", "uidx_customer_member_request"},
		{"customer_benefit_policies", "uidx_customer_policy_organization"},
		{"customer_daily_point_grant_budgets", "uidx_customer_budget_organization_date"},
		{"customer_point_entries", "uidx_customer_point_request"},
		{"customer_point_entries", "uidx_customer_point_reversal"},
		{"customer_coupon_templates", "uidx_customer_coupon_template_issuer_code"},
		{"customer_coupon_templates", "uidx_customer_coupon_template_issuer_request"},
		{"customer_coupon_grants", "uidx_customer_coupon_request"},
		{"customer_coupon_grants", "uidx_customer_coupon_issuer_request"},
		{"customer_coupon_distribution_jobs", "uidx_customer_coupon_distribution_job_request"},
	} {
		assertUniqueIndex(t, db, item.table, item.index)
	}
	for _, item := range []struct{ table, index string }{
		{"customer_coupon_distribution_jobs", "idx_customer_coupon_distribution_job_claim"},
		{"customer_members", "idx_customer_member_distribution"},
		{"customer_coupon_templates", "idx_customer_coupon_template_catchup"},
	} {
		if !db.Migrator().HasIndex(item.table, item.index) {
			t.Fatalf("missing customer lookup index %s", item.index)
		}
	}
	assertCustomerPhoneReuse(t, db)
	assertCouponIssuerKeys(t, db)
	assertCouponGrantIssuerRequest(t, db)
}

func assertCouponGrantIssuerRequest(t *testing.T, db *gorm.DB) {
	t.Helper()
	digest := "issuer-request-digest"
	first := gen.CustomerCouponGrant{ID: "grant-1", TemplateID: "template-1", MemberID: "member-1", RequestKey: "once", IssuerRequestDigest: &digest, IssuedAt: 1_800_000_000_000}
	if err := db.Create(&first).Error; err != nil {
		t.Fatal(err)
	}
	second := gen.CustomerCouponGrant{ID: "grant-2", TemplateID: "template-2", MemberID: "member-2", RequestKey: "once", IssuerRequestDigest: &digest, IssuedAt: 1_800_000_000_000}
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("duplicate issuer request accepted")
	}
}

func assertCouponIssuerKeys(t *testing.T, db *gorm.DB) {
	t.Helper()
	storeID := "store-1"
	base := gen.CustomerCouponTemplate{ID: "hq-coupon", OrganizationID: "hq", Code: "WELCOME", RequestKey: "issue-1", Title: "Welcome", IssuerScope: gen.CouponIssuerScopeHeadquarters,
		AmountFen: 1, DaysAfterActivation: 1, EffectiveAt: 1_800_000_000_000, PerMemberLimit: 1}
	if err := db.Create(&base).Error; err != nil {
		t.Fatal(err)
	}
	store := base
	store.ID = "store-coupon"
	store.IssuerScope = gen.CouponIssuerScopeStore
	store.ApplicableStoreID = &storeID
	if err := db.Create(&store).Error; err != nil {
		t.Fatalf("issuer scopes should have independent keys: %v", err)
	}
	duplicate := store
	duplicate.ID = "store-duplicate"
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("same-store coupon code accepted twice")
	}
}

func assertCustomerPhoneReuse(t *testing.T, db *gorm.DB) {
	t.Helper()
	phone := "+8613800000001"
	newMember := func(id string, value *string) gen.CustomerMember {
		return gen.CustomerMember{
			ID: id, MemberNumber: id, RequestKey: id, Phone: value,
			Status: gen.CustomerMemberStatusActive, OrganizationID: "hq",
			NoticeVersion: "v1", ProcessingBasisCode: "service",
			EvidenceReference: "evidence",
		}
	}
	first := newMember("member-1", &phone)
	if err := db.Create(&first).Error; err != nil {
		t.Fatal(err)
	}
	second := newMember("member-2", &phone)
	if err := db.Create(&second).Error; err == nil {
		t.Fatal("duplicate HQ phone was accepted")
	}
	if err := db.Model(&first).Update("phone", nil).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&second).Error; err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"cancelled-1", "cancelled-2"} {
		record := newMember(id, nil)
		record.Status = gen.CustomerMemberStatusCancelled
		if err := db.Create(&record).Error; err != nil {
			t.Fatalf("nullable cancelled phone %s: %v", id, err)
		}
	}
}
