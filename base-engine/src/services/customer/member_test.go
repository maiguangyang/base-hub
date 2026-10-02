package customer

import (
	"context"
	"strings"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/dbup"
	"base-engine/src/services/audit"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newCustomerServiceFixture(t *testing.T) (*Service, *gorm.DB, *auth.WorkspacePrincipal) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range []any{
		&gen.Organization{}, &gen.CustomerMember{}, &gen.CustomerBenefitPolicy{},
		&gen.CustomerDailyPointGrantBudget{}, &gen.CustomerPointEntry{},
		&gen.CustomerCouponTemplate{}, &gen.CustomerCouponGrant{}, &gen.CustomerCouponDistributionJob{}, &gen.AuditLog{},
	} {
		if err := db.AutoMigrate(model); err != nil {
			t.Fatal(err)
		}
		for _, index := range []string{"is_delete", "weight", "state", "deleted_by", "updated_by", "created_by"} {
			if err := db.Exec("DROP INDEX IF EXISTS `" + index + "`").Error; err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := dbup.EnsureCustomerIndexes(db); err != nil {
		t.Fatal(err)
	}
	hq := gen.Organization{ID: "hq", Code: "HQ", Name: "HQ", Type: gen.OrganizationTypeHeadquarters, Status: gen.OrganizationStatusActive}
	if err := db.Create(&hq).Error; err != nil {
		t.Fatal(err)
	}
	hqID := hq.ID
	actions := map[string]struct{}{}
	for _, action := range []string{"hqCustomer:read", "hqCustomer:create", "hqCustomer:update", "hqCustomer:cancel", "customer:read_sensitive", "hqCustomerPolicy:read", "hqCustomerPolicy:manage"} {
		actions[action] = struct{}{}
	}
	principal := &auth.WorkspacePrincipal{AccountID: "admin", SessionID: "session", WorkspaceType: auth.WorkspaceTypeHeadquarters, OrganizationID: &hqID, Permissions: actions}
	return NewService(db, audit.NewService()), db, principal
}

func customerCreateInput(key string) CreateMemberInput {
	return CreateMemberInput{Phone: "13800000001", RequestKey: key}
}

func TestCustomerMemberCreate(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	ctx := context.Background()
	member, err := service.CreateMember(ctx, principal, customerCreateInput("create-1"))
	if err != nil {
		t.Fatal(err)
	}
	if member.OrganizationID != "hq" || member.MemberNumber == "" ||
		member.Phone == nil || *member.Phone != "+8613800000001" {
		t.Fatalf("member identity = %+v", member)
	}
	assertCreatedTime(t, member.CreatedAt)
	replayed, err := service.CreateMember(ctx, principal, customerCreateInput("create-1"))
	if err != nil || replayed.ID != member.ID {
		t.Fatalf("create replay = %+v, %v", replayed, err)
	}
	var count int64
	if err := db.Model(&gen.CustomerMember{}).Count(&count).Error; err != nil || count != 1 {
		t.Fatalf("member count = %d, %v", count, err)
	}
}

func assertCreatedTime(t *testing.T, createdAt int64) {
	t.Helper()
	if createdAt <= 0 {
		t.Fatalf("missing created time: %d", createdAt)
	}
}

func TestCustomerMemberDuplicatePhone(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	if _, err := service.CreateMember(context.Background(), principal, customerCreateInput("create-1")); err != nil {
		t.Fatal(err)
	}
	_, err := service.CreateMember(context.Background(), principal, customerCreateInput("create-2"))
	if auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("duplicate phone code = %v", err)
	}
}

func TestCustomerMemberAllowsCreationWithoutNoticeEvidence(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	input := CreateMemberInput{Phone: "13800000001", RequestKey: "phone-only"}
	member, err := service.CreateMember(context.Background(), principal, input)
	if err != nil {
		t.Fatal(err)
	}
	if member.NoticeVersion != "" || member.ProcessingBasisCode != "" || member.EvidenceReference != "" {
		t.Fatalf("unexpected notice evidence: %+v", member)
	}
	var record gen.AuditLog
	if err := db.Where("action = ? AND resource_id = ?", "hqCustomer:create", member.ID).First(&record).Error; err != nil {
		t.Fatal(err)
	}
	if record.MetadataJSON != nil && strings.Contains(*record.MetadataJSON, "evidenceReference") {
		t.Fatalf("fabricated evidence in audit: %s", *record.MetadataJSON)
	}
	replay, err := service.CreateMember(context.Background(), principal, input)
	if err != nil || replay.ID != member.ID {
		t.Fatalf("phone-only replay = %+v, %v", replay, err)
	}
}

func TestCustomerMemberCreateRequiresElevenMainlandDigits(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	for _, phone := range []string{"1380000000", "138000000000", "12800000000", "1380000000a", "+8613800000001", "138-0000 0001"} {
		_, err := service.CreateMember(context.Background(), principal, CreateMemberInput{Phone: phone, RequestKey: "invalid-" + phone})
		if auth.ErrorCode(err) != auth.CodeValidationFailed {
			t.Errorf("invalid member phone %q accepted: %v", phone, err)
		}
	}
}

func TestCustomerMemberRequiresHQ(t *testing.T) {
	service, _, principal := newCustomerServiceFixture(t)
	principal.WorkspaceType = auth.WorkspaceTypeFranchise
	if _, err := service.CreateMember(context.Background(), principal, customerCreateInput("create-1")); auth.ErrorCode(err) != auth.CodeWorkspaceForbidden {
		t.Fatalf("franchise creation code = %v", err)
	}
}

func TestCustomerMemberCancellationRetainsRights(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	ctx := context.Background()
	member, err := service.CreateMember(ctx, principal, customerCreateInput("create-1"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetMemberStatus(ctx, principal, member.ID, gen.CustomerMemberStatusCancelled); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("direct cancellation code = %v", err)
	}
	if _, err := service.RequestCancellation(ctx, principal, member.ID, "identity-check-1", "MEMBERSHIP_END"); err != nil {
		t.Fatal(err)
	}
	assertCancellationAuditBasis(t, db)
	assertPointsBlockCancellation(t, service, db, principal, member.ID)
	assertCouponBlocksCancellation(t, service, db, principal, member.ID)
	if _, err := service.CompleteCancellation(ctx, principal, member.ID, "settlement-1"); err != nil {
		t.Fatal(err)
	}
	var cancelled gen.CustomerMember
	if err := db.First(&cancelled, "id = ?", member.ID).Error; err != nil || cancelled.Phone != nil {
		t.Fatalf("cancelled member phone = %v, %v", cancelled.Phone, err)
	}
	rejoined, err := service.CreateMember(ctx, principal, customerCreateInput("create-2"))
	if err != nil || rejoined.ID == member.ID {
		t.Fatalf("rejoined member = %+v, %v", rejoined, err)
	}
}

func TestCustomerMemberCancellationRejectsFrozenPointMismatch(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	ctx := context.Background()
	member, err := service.CreateMember(ctx, principal, customerCreateInput("frozen-cancel"))
	if err != nil {
		t.Fatal(err)
	}
	entry := gen.CustomerPointEntry{
		ID: "unsettled-point", MemberID: member.ID, SourceOrganizationID: "hq",
		Delta: 10, RequestKey: "unsettled-point",
		Source: gen.CustomerPointSourceHqManual, OperationKind: gen.CustomerPointOperationKindGrant,
		ReasonCode: gen.CustomerPointReasonCodeReward,
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatal(err)
	}
	if err := service.ReconcileMemberPoints(ctx, member.ID); auth.ErrorCode(err) != auth.CodeCustomerPointsMismatch {
		t.Fatalf("reconciliation = %v", err)
	}
	if _, err := service.RequestCancellation(ctx, principal, member.ID, "identity-check", "MEMBERSHIP_END"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.CompleteCancellation(ctx, principal, member.ID, "settlement-claim"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("frozen ledger mismatch must block cancellation: %v", err)
	}
}

func TestCustomerMemberCancellationRejectsUnreconciledPointMismatch(t *testing.T) {
	service, db, principal := newCustomerServiceFixture(t)
	ctx := context.Background()
	member, err := service.CreateMember(ctx, principal, customerCreateInput("unreconciled-cancel"))
	if err != nil {
		t.Fatal(err)
	}
	entry := gen.CustomerPointEntry{
		ID: "unreconciled-point", MemberID: member.ID, SourceOrganizationID: "hq",
		Delta: 10, RequestKey: "unreconciled-point",
		Source: gen.CustomerPointSourceHqManual, OperationKind: gen.CustomerPointOperationKindGrant,
		ReasonCode: gen.CustomerPointReasonCodeReward,
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.RequestCancellation(ctx, principal, member.ID, "identity-check", "MEMBERSHIP_END"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.CompleteCancellation(ctx, principal, member.ID, "settlement-claim"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("unreconciled ledger mismatch must block cancellation: %v", err)
	}
}

func assertPointsBlockCancellation(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal, memberID string) {
	t.Helper()
	if err := db.Model(&gen.CustomerMember{}).Where("id = ?", memberID).Update("points_balance", 10).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.CompleteCancellation(context.Background(), principal, memberID, "settlement-1"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("positive points must block cancellation: %v", err)
	}
	if err := db.Model(&gen.CustomerMember{}).Where("id = ?", memberID).Update("points_balance", 0).Error; err != nil {
		t.Fatal(err)
	}
}

func assertCouponBlocksCancellation(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal, memberID string) {
	t.Helper()
	template := gen.CustomerCouponTemplate{
		ID: "template-1", Code: "C1", Title: "Coupon", OrganizationID: "hq",
		RequestKey: "template-1", AmountFen: 100, DaysAfterActivation: 30,
		EffectiveAt: 1_800_000_000_000, PerMemberLimit: 1, TotalIssueLimit: 1,
	}
	if err := db.Create(&template).Error; err != nil {
		t.Fatal(err)
	}
	grant := gen.CustomerCouponGrant{
		ID: "grant-1", MemberID: memberID, TemplateID: template.ID,
		Status:    gen.CustomerCouponGrantStatusPendingActivation,
		AmountFen: 100, DaysAfterActivation: 30, RequestKey: "grant-1", IssuedAt: 1_800_000_000_000,
	}
	if err := db.Create(&grant).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := service.CompleteCancellation(context.Background(), principal, memberID, "settlement-1"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("pending coupon must block cancellation: %v", err)
	}
	if err := db.Model(&grant).Update("status", gen.CustomerCouponGrantStatusRevoked).Error; err != nil {
		t.Fatal(err)
	}
}
