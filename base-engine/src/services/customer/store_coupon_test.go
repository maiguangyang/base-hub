package customer

import (
	"context"
	"strings"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func storeCouponFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal, *auth.WorkspacePrincipal, *gen.CustomerMember) {
	service, db, hq, member := couponFixture(t)
	if err := db.AutoMigrate(&gen.Store{}); err != nil {
		t.Fatal(err)
	}
	for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
		db.Exec("DROP INDEX IF EXISTS `" + index + "`")
	}
	org := gen.Organization{ID: "franchise", Code: "FRANCHISE", Name: "Franchise", Type: gen.OrganizationTypeFranchise, Status: gen.OrganizationStatusActive}
	store := gen.Store{ID: "store", Code: "S", Name: "Store", OrganizationID: org.ID, Lifecycle: gen.StoreLifecycleActive}
	for _, item := range []any{&org, &store} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
	orgID := org.ID
	principal := &auth.WorkspacePrincipal{AccountID: "operator", SessionID: "session", WorkspaceType: auth.WorkspaceTypeFranchise, OrganizationID: &orgID,
		Permissions: map[string]struct{}{"franchiseCoupon:read": {}, "franchiseCoupon:manage": {}, "franchiseCoupon:grant": {}, "franchiseCoupon:revoke": {}}, StoreIDs: map[string]struct{}{store.ID: {}}}
	return service, db, hq, principal, member
}

func TestStoreCouponEffectiveTimeGate(t *testing.T) {
	service, db, _, principal, member := storeCouponFixture(t)
	ctx := context.Background()
	effectiveAt := int64(1_800_000_000_000)
	distributionEndsAt := effectiveAt + 1_000
	input := validCouponInput()
	input.Code, input.RequestKey, input.EffectiveAt = "STORE-FUTURE", "store-future", &effectiveAt
	input.DistributionEndsAt, input.PerMemberLimit = &distributionEndsAt, 2
	template, err := service.CreateStoreCouponTemplate(ctx, principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return time.UnixMilli(effectiveAt - 1) }
	assertStoreCouponGrantConflict(t, service, principal, template.ID, member.ID, "before")
	assertCouponCounts(t, db, template.ID, 0, 0)
	service.now = func() time.Time { return time.UnixMilli(effectiveAt) }
	grant, err := service.GrantStoreCoupon(ctx, principal, "store", template.ID, member.ID, "exact")
	if err != nil {
		t.Fatal(err)
	}
	if grant.IssuedAt != effectiveAt {
		t.Fatalf("issued_at = %d, want %d", grant.IssuedAt, effectiveAt)
	}
	service.now = func() time.Time { return time.UnixMilli(distributionEndsAt) }
	assertStoreCouponGrantConflict(t, service, principal, template.ID, member.ID, "at-cutoff")
	replay, err := service.GrantStoreCoupon(ctx, principal, "store", template.ID, member.ID, "exact")
	if err != nil || replay.ID != grant.ID {
		t.Fatalf("expired store template replay = %+v, %v", replay, err)
	}
}

func assertStoreCouponGrantConflict(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, templateID, memberID, key string) {
	t.Helper()
	if _, err := service.GrantStoreCoupon(context.Background(), principal, "store", templateID, memberID, key); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("store grant %s = %v", key, err)
	}
}

func TestStoreCouponSharesMemberButSeparatesIssuer(t *testing.T) {
	service, _, hq, principal, member := storeCouponFixture(t)
	input := validCouponInput()
	input.RequestKey = "store-template"
	input.Code = "STORE-WELCOME"
	template, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	assertStoreTemplateAndHqSeparation(t, service, hq, template)
	assertStoreGrantReplayAndScope(t, service, hq, principal, member, template)
}

func TestStoreCouponTemplateReplayIntentMatchesOriginalRequest(t *testing.T) {
	service, _, _, principal, _ := storeCouponFixture(t)
	input := validCouponInput()
	input.Code, input.RequestKey, input.EffectiveAt = "STORE-REPLAY", "store-template-replay", nil
	template, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	replay, err := findStoreCouponTemplateByRequest(service.db, template.OrganizationID, "store", input.RequestKey)
	if err != nil || replay == nil || replay.ID != template.ID || !sameTemplateIntent(replay, input) {
		t.Fatalf("matching store template replay = %+v, %v", replay, err)
	}
	changed := input
	changed.Title = "Changed"
	if sameTemplateIntent(replay, changed) {
		t.Fatal("store template replay accepted changed intent")
	}
}

func TestStoreCouponAutomaticGrantUsesStoreAuditAttribution(t *testing.T) {
	service, db, _, principal, _ := storeCouponFixture(t)
	input := validCouponInput()
	input.Code, input.RequestKey = "AUTO-STORE", "auto-store-template"
	template, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	hqID := template.OrganizationID
	other := gen.CustomerMember{ID: "store-auto-member", MemberNumber: "CM-STORE-AUTO", RequestKey: "store-auto-member",
		Phone: stringPointer("13800000003"), Status: gen.CustomerMemberStatusActive, OrganizationID: hqID,
		CreatedAt: time.Now().UnixMilli()}
	if err := db.Create(&other).Error; err != nil {
		t.Fatal(err)
	}
	storeID := "store"
	var grant *gen.CustomerCouponGrant
	var created bool
	err = db.Transaction(func(tx *gorm.DB) error {
		grant, created, err = service.createCouponGrantInTransaction(tx, couponGrantInput{
			TemplateID: template.ID, MemberID: other.ID, RequestKey: autoCouponRequestKey(template.ID, other.ID),
			ExpectedStoreID: &storeID, Automatic: true, NowMillis: time.Now().UnixMilli(),
		})
		return err
	})
	if err != nil || !created {
		t.Fatalf("automatic store grant = %+v created:%t err:%v", grant, created, err)
	}
	assertAutomaticCouponAudit(t, db, grant.ID, "franchiseCoupon:grant", hqID, &storeID)
	if !strings.HasPrefix(grant.RequestKey, "AUTO:") {
		t.Fatalf("automatic request key = %q", grant.RequestKey)
	}
}

func stringPointer(value string) *string { return &value }

func TestStoreCouponRequestKeyIsBoundToIssuerAndPayload(t *testing.T) {
	service, _, hq, principal, member := storeCouponFixture(t)
	input := validCouponInput()
	input.Code, input.RequestKey = "STORE-ONE", "template-one"
	first, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	input.Code, input.RequestKey = "STORE-TWO", "template-two"
	second, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantStoreCoupon(context.Background(), principal, "store", first.ID, member.ID, "issue-once"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantStoreCoupon(context.Background(), principal, "store", second.ID, member.ID, "issue-once"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("cross-template replay: %v", err)
	}
	otherInput := customerCreateInput("member-two")
	otherInput.Phone = "13800000002"
	other, err := service.CreateMember(context.Background(), hq, otherInput)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.GrantStoreCoupon(context.Background(), principal, "store", first.ID, other.ID, "issue-once"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("cross-member replay: %v", err)
	}
}

func TestStoreCouponManagementDoesNotCrossStores(t *testing.T) {
	service, db, _, principal, member := storeCouponFixture(t)
	otherStore := gen.Store{ID: "other-store", Code: "OTHER", Name: "Other Store", OrganizationID: "franchise", Lifecycle: gen.StoreLifecycleActive}
	if err := db.Create(&otherStore).Error; err != nil {
		t.Fatal(err)
	}
	principal.StoreIDs[otherStore.ID] = struct{}{}
	input := validCouponInput()
	input.RequestKey, input.Code = "store-one-template", "STORE-ONE"
	template, err := service.CreateStoreCouponTemplate(context.Background(), principal, "store", input)
	if err != nil {
		t.Fatal(err)
	}
	grant, err := service.GrantStoreCoupon(context.Background(), principal, "store", template.ID, member.ID, "store-one-grant")
	if err != nil {
		t.Fatal(err)
	}

	otherTemplates, err := service.ListStoreCouponTemplates(context.Background(), principal, otherStore.ID, 1, 20)
	if err != nil || otherTemplates.Total != 0 {
		t.Fatalf("other store templates = %+v, %v", otherTemplates, err)
	}
	if _, err := service.GrantStoreCoupon(context.Background(), principal, otherStore.ID, template.ID, member.ID, "cross-store-grant"); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("other store granted coupon: %v", err)
	}
	if _, err := service.SetStoreCouponTemplateEnabled(context.Background(), principal, otherStore.ID, template.ID, false); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("other store changed template: %v", err)
	}
	if _, err := service.ListStoreCouponGrants(context.Background(), principal, otherStore.ID, template.ID, 1, 20); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("other store listed grants: %v", err)
	}
	if _, err := service.RevokeStoreCoupon(context.Background(), principal, otherStore.ID, grant.ID, "CROSS_STORE"); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("other store revoked grant: %v", err)
	}
}

func assertStoreTemplateAndHqSeparation(t *testing.T, service *Service, hq *auth.WorkspacePrincipal, template *gen.CustomerCouponTemplate) {
	t.Helper()
	if template.IssuerScope != gen.CouponIssuerScopeStore || template.ApplicableStoreID == nil || *template.ApplicableStoreID != "store" || template.OrganizationID != "hq" {
		t.Fatalf("template = %+v", template)
	}
	hqTemplates, err := service.ListCouponTemplates(context.Background(), hq, 1, 20)
	if err != nil || hqTemplates.Total != 0 {
		t.Fatalf("HQ templates = %+v, %v", hqTemplates, err)
	}
}

func assertStoreGrantReplayAndScope(t *testing.T, service *Service, hq, principal *auth.WorkspacePrincipal, member *gen.CustomerMember, template *gen.CustomerCouponTemplate) {
	t.Helper()
	grant, err := service.GrantStoreCoupon(context.Background(), principal, "store", template.ID, member.ID, "store-grant")
	if err != nil {
		t.Fatal(err)
	}
	if grant.Status != gen.CustomerCouponGrantStatusPendingActivation {
		t.Fatalf("grant = %+v", grant)
	}
	replay, err := service.GrantStoreCoupon(context.Background(), principal, "store", template.ID, member.ID, "store-grant")
	if err != nil || replay.ID != grant.ID {
		t.Fatalf("replay = %+v, %v", replay, err)
	}
	if _, err := service.GrantCoupon(context.Background(), hq, template.ID, member.ID, "hq-grant"); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("HQ granted store template: %v", err)
	}
	principal.StoreIDs = map[string]struct{}{}
	if _, err := service.GrantStoreCoupon(context.Background(), principal, "store", template.ID, member.ID, "other"); auth.ErrorCode(err) != auth.CodeStoreScopeDenied {
		t.Fatalf("scope: %v", err)
	}
}

func TestStoreCouponMemberLookupRequiresExactIdentifierAndLimitsRequests(t *testing.T) {
	service, _, _, principal, member := storeCouponFixture(t)
	if _, err := service.LookupStoreCouponMember(context.Background(), principal, "store", "1380000"); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("prefix lookup: %v", err)
	}
	found, err := service.LookupStoreCouponMember(context.Background(), principal, "store", member.ID)
	if err != nil || found == nil || found.ID != member.ID || found.PhoneMasked == "" || found.PhoneMasked == *member.Phone {
		t.Fatalf("lookup = %+v, %v", found, err)
	}
	for i := 0; i < 19; i++ {
		if _, err := service.LookupStoreCouponMember(context.Background(), principal, "store", member.ID); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := service.LookupStoreCouponMember(context.Background(), principal, "store", member.ID); auth.ErrorCode(err) != auth.CodeRateLimited {
		t.Fatalf("lookup limit: %v", err)
	}
}
