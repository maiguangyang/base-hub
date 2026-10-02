package customer

import (
	"context"
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) ReversePoints(ctx context.Context, principal *auth.WorkspacePrincipal, entryID string, reason gen.CustomerPointReasonCode, note, key string) (*PointMutationResult, error) {
	hqID, err := headquartersID(principal, "hqCustomerPoints:reverse", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	if !reason.IsValid() || strings.TrimSpace(note) == "" || len(note) > 512 ||
		strings.TrimSpace(key) == "" || len(key) > 128 {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	return s.runPointWrite(ctx,
		func(tx *gorm.DB) (*gen.CustomerPointEntry, error) {
			return s.reversePointsInTransaction(tx, principal, hqID, entryID, reason, note, key)
		},
		func() (*gen.CustomerPointEntry, error) {
			return s.replayReverseAfterConflict(ctx, hqID, entryID, reason, note, key)
		},
	)
}

func (s *Service) replayReverseAfterConflict(ctx context.Context, hqID, entryID string, reason gen.CustomerPointReasonCode, note, key string) (*gen.CustomerPointEntry, error) {
	entry, err := findPointByRequest(s.db.WithContext(ctx), hqID, gen.CustomerPointOperationKindReverse, key)
	if err != nil {
		return nil, err
	}
	if entry == nil || !sameReverseIntent(entry, entryID, reason, note) {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return entry, nil
}

func sameReverseIntent(entry *gen.CustomerPointEntry, entryID string, reason gen.CustomerPointReasonCode, note string) bool {
	return entry.ReversesID != nil && *entry.ReversesID == entryID &&
		entry.ReasonCode == reason && entry.Note != nil && *entry.Note == note
}

func (s *Service) reversePointsInTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, entryID string, reason gen.CustomerPointReasonCode, note, key string) (*gen.CustomerPointEntry, error) {
	replayed, err := findPointByRequest(tx, hqID, gen.CustomerPointOperationKindReverse, key)
	if err != nil {
		return nil, err
	}
	if replayed != nil {
		if !sameReverseIntent(replayed, entryID, reason, note) {
			return nil, auth.NewError(auth.CodeConflict)
		}
		return replayed, nil
	}
	original, member, err := loadReversiblePoint(tx, hqID, entryID)
	if err != nil {
		return nil, err
	}
	entry := &gen.CustomerPointEntry{
		ID: uuid.Must(uuid.NewV4()).String(), Delta: -original.Delta,
		CreatedAt: s.now().UnixMilli(),
		Source:    gen.CustomerPointSourceHqManual, OperationKind: gen.CustomerPointOperationKindReverse,
		ReasonCode: reason, Note: &note, RequestKey: key, MemberID: member.ID,
		SourceOrganizationID: hqID, ReversesID: &original.ID,
	}
	if err := tx.Create(entry).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(member).Update("points_balance", member.PointsBalance-original.Delta).Error; err != nil {
		return nil, err
	}
	return entry, s.auditPoint(tx, principal, entry)
}

func loadReversiblePoint(tx *gorm.DB, hqID, entryID string) (*gen.CustomerPointEntry, *gen.CustomerMember, error) {
	var original gen.CustomerPointEntry
	err := tx.Where("id = ? AND source_organization_id = ?", entryID, hqID).First(&original).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return nil, nil, err
	}
	if original.OperationKind != gen.CustomerPointOperationKindGrant || original.Source != gen.CustomerPointSourceHqManual {
		return nil, nil, auth.NewError(auth.CodeConflict)
	}
	var member gen.CustomerMember
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND organization_id = ?", original.MemberID, hqID).First(&member).Error
	if err != nil {
		return nil, nil, err
	}
	if member.PointsFrozen || member.PointsBalance < original.Delta {
		return nil, nil, auth.NewError(auth.CodeConflict)
	}
	var count int64
	if err := tx.Model(&gen.CustomerPointEntry{}).Where("reverses_id = ?", original.ID).Count(&count).Error; err != nil {
		return nil, nil, err
	}
	if count != 0 {
		return nil, nil, auth.NewError(auth.CodeConflict)
	}
	return &original, &member, nil
}
