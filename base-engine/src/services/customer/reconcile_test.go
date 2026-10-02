package customer

import (
	"context"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func TestCustomerPointReconciliation(t *testing.T) {
	service, db, principal, member := pointsFixture(t)
	configurePoints(t, service, principal, 100, 100)
	ctx := context.Background()
	grant, err := service.GrantPoints(ctx, principal, member.ID, 50, gen.CustomerPointReasonCodeReward, "gift", "gift-1")
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(member).Update("points_balance", 70).Error; err != nil {
		t.Fatal(err)
	}
	assertCustomerPointFrozen(t, service, db, principal, member.ID)
	assertFrozenReversalBlocked(t, service, principal, grant.ID)
	correction, err := service.CorrectPointMismatch(ctx, principal, member.ID, "investigation-1", "correct-1")
	if err != nil || correction.Delta != 20 {
		t.Fatalf("correction = %+v, %v", correction, err)
	}
	replay, err := service.CorrectPointMismatch(ctx, principal, member.ID, "investigation-1", "correct-1")
	if err != nil || replay.ID != correction.ID {
		t.Fatalf("correction replay = %+v, %v", replay, err)
	}
	assertChangedCorrectionEvidence(t, service, principal, member.ID)
	var reconciled gen.CustomerMember
	if err := db.First(&reconciled, "id = ?", member.ID).Error; err != nil || reconciled.PointsFrozen || reconciled.PointsBalance != 70 {
		t.Fatalf("reconciled = %+v, %v", reconciled, err)
	}
}

func assertFrozenReversalBlocked(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, entryID string) {
	t.Helper()
	if _, err := service.ReversePoints(context.Background(), principal, entryID, gen.CustomerPointReasonCodeCorrection, "undo", "reverse-frozen"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("frozen point reversal = %v", err)
	}
}

func assertChangedCorrectionEvidence(t *testing.T, service *Service, principal *auth.WorkspacePrincipal, memberID string) {
	t.Helper()
	if _, err := service.CorrectPointMismatch(context.Background(), principal, memberID, "different-evidence", "correct-1"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("changed correction evidence = %v", err)
	}
}

func assertCustomerPointFrozen(t *testing.T, service *Service, db *gorm.DB, principal *auth.WorkspacePrincipal, memberID string) {
	t.Helper()
	ctx := context.Background()
	if err := service.ReconcileMemberPoints(ctx, memberID); auth.ErrorCode(err) != auth.CodeCustomerPointsMismatch {
		t.Fatalf("mismatch code = %v", err)
	}
	var frozen gen.CustomerMember
	if err := db.First(&frozen, "id = ?", memberID).Error; err != nil || !frozen.PointsFrozen {
		t.Fatalf("frozen = %+v, %v", frozen, err)
	}
	if _, err := service.GrantPoints(ctx, principal, memberID, 1, gen.CustomerPointReasonCodeReward, "gift", "blocked"); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("frozen gift = %v", err)
	}
}

func TestCustomerPointCorrectionRequiresEvidence(t *testing.T) {
	service, _, principal, member := pointsFixture(t)
	if _, err := service.CorrectPointMismatch(context.Background(), principal, member.ID, "", "correct-1"); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("blank evidence = %v", err)
	}
	delete(principal.Permissions, "hqCustomerPoints:correct")
	if _, err := service.CorrectPointMismatch(context.Background(), principal, member.ID, "proof", "correct-1"); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("missing correct action = %v", err)
	}
}

func TestCustomerPointReconciliationScansAndAlertsOnlyAffectedMember(t *testing.T) {
	service, db, brokenID, healthyID := newReconciliationScanFixture(t)
	var alerts []string
	for range 2 {
		if err := service.ReconcileAllMemberPoints(context.Background(), func(id string, cause error) {
			if auth.ErrorCode(cause) != auth.CodeCustomerPointsMismatch {
				t.Errorf("unexpected reconciliation error for %s: %v", id, cause)
			}
			alerts = append(alerts, id)
		}); err != nil {
			t.Fatal(err)
		}
	}
	if len(alerts) != 2 || alerts[0] != brokenID || alerts[1] != brokenID {
		t.Fatalf("alerts = %v", alerts)
	}
	assertMemberPointsFrozen(t, db, brokenID, true)
	assertMemberPointsFrozen(t, db, healthyID, false)
	assertSingleReconciliationAudit(t, db, brokenID)
}

func newReconciliationScanFixture(t *testing.T) (*Service, *gorm.DB, string, string) {
	t.Helper()
	service, db, principal := newCustomerServiceFixture(t)
	broken, err := service.CreateMember(context.Background(), principal, customerCreateInput("broken-member"))
	if err != nil {
		t.Fatal(err)
	}
	healthyInput := customerCreateInput("healthy-member")
	healthyInput.Phone = "13900000002"
	healthy, err := service.CreateMember(context.Background(), principal, healthyInput)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Model(broken).Update("points_balance", 25).Error; err != nil {
		t.Fatal(err)
	}
	return service, db, broken.ID, healthy.ID
}

func assertMemberPointsFrozen(t *testing.T, db *gorm.DB, memberID string, frozen bool) {
	t.Helper()
	var member gen.CustomerMember
	if err := db.First(&member, "id = ?", memberID).Error; err != nil || member.PointsFrozen != frozen {
		t.Fatalf("member %s frozen = %v, %v", memberID, member.PointsFrozen, err)
	}
}

func assertSingleReconciliationAudit(t *testing.T, db *gorm.DB, memberID string) {
	t.Helper()
	var audits int64
	if err := db.Model(&gen.AuditLog{}).Where("action = ? AND resource_id = ?", "hqCustomerPoints:reconcile", memberID).Count(&audits).Error; err != nil || audits != 1 {
		t.Fatalf("duplicate reconciliation audit count = %d, %v", audits, err)
	}
}
