package customer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func pointLedgerSum(tx *gorm.DB, memberID string) (int64, error) {
	var sum int64
	err := tx.Model(&gen.CustomerPointEntry{}).
		Where("member_id = ?", memberID).
		Select("COALESCE(SUM(delta), 0)").Scan(&sum).Error
	return sum, err
}

// ReconcileAllMemberPoints scans every member in bounded pages. A failed member
// does not prevent the remaining members from being checked.
func (s *Service) ReconcileAllMemberPoints(ctx context.Context, report func(string, error)) error {
	const pageSize = 100
	lastID := ""
	var failures error
	for {
		var ids []string
		err := s.db.WithContext(ctx).Model(&gen.CustomerMember{}).
			Where("id > ?", lastID).Order("id").Limit(pageSize).Pluck("id", &ids).Error
		if err != nil {
			return errors.Join(failures, err)
		}
		failures = errors.Join(failures, s.reconcileMemberIDs(ctx, ids, report))
		if err := ctx.Err(); err != nil {
			return errors.Join(failures, err)
		}
		if len(ids) < pageSize {
			return failures
		}
		lastID = ids[len(ids)-1]
	}
}

func (s *Service) reconcileMemberIDs(ctx context.Context, ids []string, report func(string, error)) error {
	var failures error
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return errors.Join(failures, err)
		}
		if err := s.ReconcileMemberPoints(ctx, id); err != nil {
			if report != nil {
				report(id, err)
			}
			if auth.ErrorCode(err) != auth.CodeCustomerPointsMismatch {
				failures = errors.Join(failures, fmt.Errorf("member %s: %w", id, err))
			}
		}
	}
	return failures
}

func (s *Service) ReconcileMemberPoints(ctx context.Context, memberID string) error {
	mismatched := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var member gen.CustomerMember
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&member, "id = ?", memberID).Error; err != nil {
			return err
		}
		sum, err := pointLedgerSum(tx, memberID)
		if err != nil {
			return err
		}
		if sum == member.PointsBalance {
			return nil
		}
		mismatched = true
		if member.PointsFrozen {
			return nil
		}
		if err := tx.Model(&member).Update("points_frozen", true).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{
			OrganizationID: &member.OrganizationID, Action: "hqCustomerPoints:reconcile",
			ResourceType: "customerMember", ResourceID: member.ID, ResultCode: "MISMATCH",
			Metadata: audit.Metadata{ExpectedPoints: member.PointsBalance, LedgerPoints: sum},
		})
	})
	if err != nil {
		return err
	}
	if mismatched {
		return auth.NewError(auth.CodeCustomerPointsMismatch)
	}
	return nil
}

func (s *Service) CorrectPointMismatch(ctx context.Context, principal *auth.WorkspacePrincipal, memberID, evidenceReference, key string) (*PointMutationResult, error) {
	hqID, err := headquartersID(principal, "hqCustomerPoints:correct", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(evidenceReference) == "" || len(evidenceReference) > 128 ||
		strings.TrimSpace(key) == "" || len(key) > 128 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	return s.runPointWrite(ctx,
		func(tx *gorm.DB) (*gen.CustomerPointEntry, error) {
			return s.correctPointInTransaction(tx, principal, hqID, memberID, evidenceReference, key)
		},
		func() (*gen.CustomerPointEntry, error) {
			return s.replayCorrectionAfterConflict(ctx, hqID, memberID, evidenceReference, key)
		},
	)
}

func (s *Service) replayCorrectionAfterConflict(ctx context.Context, hqID, memberID, evidenceReference, key string) (*gen.CustomerPointEntry, error) {
	entry, err := findPointByRequest(s.db.WithContext(ctx), hqID, gen.CustomerPointOperationKindCorrect, key)
	if err != nil {
		return nil, err
	}
	if entry == nil || !sameCorrectionIntent(entry, memberID, evidenceReference) {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return entry, nil
}

func sameCorrectionIntent(entry *gen.CustomerPointEntry, memberID, evidenceReference string) bool {
	return entry.MemberID == memberID && entry.Note != nil && *entry.Note == evidenceReference
}

func (s *Service) correctPointInTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, memberID, evidenceReference, key string) (*gen.CustomerPointEntry, error) {
	replay, err := findPointByRequest(tx, hqID, gen.CustomerPointOperationKindCorrect, key)
	if err != nil {
		return nil, err
	}
	if replay != nil {
		if !sameCorrectionIntent(replay, memberID, evidenceReference) {
			return nil, auth.NewError(auth.CodeConflict)
		}
		return replay, nil
	}
	member, delta, err := correctionDifference(tx, hqID, memberID)
	if err != nil {
		return nil, err
	}
	entry := &gen.CustomerPointEntry{
		ID: uuid.Must(uuid.NewV4()).String(), Delta: delta,
		CreatedAt: s.now().UnixMilli(),
		Source:    gen.CustomerPointSourceHqManual, OperationKind: gen.CustomerPointOperationKindCorrect,
		ReasonCode: gen.CustomerPointReasonCodeCorrection, RequestKey: key, Note: &evidenceReference,
		MemberID: memberID, SourceOrganizationID: hqID,
	}
	if err := tx.Create(entry).Error; err != nil {
		return nil, err
	}
	after, err := pointLedgerSum(tx, memberID)
	if err != nil {
		return nil, err
	}
	if after != member.PointsBalance {
		return nil, auth.NewError(auth.CodeCustomerPointsMismatch)
	}
	if err := tx.Model(&member).Update("points_frozen", false).Error; err != nil {
		return nil, err
	}
	return entry, s.auditCorrection(tx, principal, entry, evidenceReference)
}

func correctionDifference(tx *gorm.DB, hqID, memberID string) (*gen.CustomerMember, int64, error) {
	var member gen.CustomerMember
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND organization_id = ?", memberID, hqID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return nil, 0, err
	}
	if !member.PointsFrozen || member.PointsBalance < 0 {
		return nil, 0, auth.NewError(auth.CodeConflict)
	}
	sum, err := pointLedgerSum(tx, memberID)
	if err != nil {
		return nil, 0, err
	}
	delta := member.PointsBalance - sum
	if delta == 0 || delta < -math.MaxInt32 || delta > math.MaxInt32 {
		return nil, 0, auth.NewError(auth.CodeConflict)
	}
	return &member, delta, nil
}

func (s *Service) auditCorrection(tx *gorm.DB, principal *auth.WorkspacePrincipal, entry *gen.CustomerPointEntry, evidence string) error {
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: principal.OrganizationID, Action: "hqCustomerPoints:correct",
		ResourceType: "customerPointEntry", ResourceID: entry.ID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{EvidenceReference: evidence}),
	})
}
