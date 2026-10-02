package customer

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var shanghaiLocation = time.FixedZone("CST", 8*3600)

func validPointRequest(points int64, reason gen.CustomerPointReasonCode, note, key string) bool {
	return points > 0 && points <= math.MaxInt32 && reason.IsValid() &&
		strings.TrimSpace(note) != "" && len(note) <= 512 &&
		strings.TrimSpace(key) != "" && len(key) <= 128
}

func (s *Service) GrantPoints(ctx context.Context, principal *auth.WorkspacePrincipal, memberID string, points int64, reason gen.CustomerPointReasonCode, note, key string) (*PointMutationResult, error) {
	hqID, err := headquartersID(principal, "hqCustomerPoints:grant", authorization.AccessCreate)
	if err != nil {
		return nil, err
	}
	if !validPointRequest(points, reason, note, key) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var result *PointMutationResult
	for attempt := 0; attempt < 3; attempt++ {
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			entry, err := s.grantPointsInTransaction(tx, principal, hqID, memberID, points, reason, note, key)
			if err != nil {
				return err
			}
			result, err = pointMutationResult(tx, entry)
			return err
		})
		if err == nil {
			return result, nil
		}
		if mysqlErrorNumber(err) == 1062 {
			entry, replayErr := s.replayGrantAfterConflict(ctx, hqID, memberID, points, reason, note, key)
			if replayErr != nil {
				return nil, replayErr
			}
			return pointMutationResult(s.db.WithContext(ctx), entry)
		}
		if mysqlErrorNumber(err) != 1213 && mysqlErrorNumber(err) != 1205 {
			return nil, err
		}
	}
	return nil, auth.NewError(auth.CodeConflict)
}

func (s *Service) replayGrantAfterConflict(ctx context.Context, hqID, memberID string, points int64, reason gen.CustomerPointReasonCode, note, key string) (*gen.CustomerPointEntry, error) {
	entry, err := findPointByRequest(s.db.WithContext(ctx), hqID, gen.CustomerPointOperationKindGrant, key)
	if err != nil {
		return nil, err
	}
	if entry == nil || entry.MemberID != memberID || entry.Delta != points ||
		entry.ReasonCode != reason || entry.Note == nil || *entry.Note != note {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return entry, nil
}

func (s *Service) grantPointsInTransaction(tx *gorm.DB, principal *auth.WorkspacePrincipal, hqID, memberID string, points int64, reason gen.CustomerPointReasonCode, note, key string) (*gen.CustomerPointEntry, error) {
	replayed, err := findPointByRequest(tx, hqID, gen.CustomerPointOperationKindGrant, key)
	if err != nil {
		return nil, err
	}
	if replayed != nil {
		if replayed.MemberID != memberID || replayed.Delta != points || replayed.ReasonCode != reason || replayed.Note == nil || *replayed.Note != note {
			return nil, auth.NewError(auth.CodeConflict)
		}
		return replayed, nil
	}
	budget, member, err := s.prepareGrant(tx, hqID, memberID, points)
	if err != nil {
		return nil, err
	}
	entry := &gen.CustomerPointEntry{
		ID: uuid.Must(uuid.NewV4()).String(), Delta: points,
		CreatedAt: s.now().UnixMilli(),
		Source:    gen.CustomerPointSourceHqManual, OperationKind: gen.CustomerPointOperationKindGrant,
		ReasonCode: reason, Note: &note, RequestKey: key,
		MemberID: memberID, SourceOrganizationID: hqID,
	}
	return entry, s.persistGrant(tx, principal, budget, member, entry)
}

func (s *Service) prepareGrant(tx *gorm.DB, hqID, memberID string, points int64) (*gen.CustomerDailyPointGrantBudget, *gen.CustomerMember, error) {
	policy, err := lockedGrantPolicy(tx, hqID, points)
	if err != nil {
		return nil, nil, err
	}
	budget, err := s.lockDailyBudget(tx, hqID)
	if err != nil {
		return nil, nil, err
	}
	if budget.UsedPoints > policy.ManualGrantMaxDaily-points {
		return nil, nil, auth.NewError(auth.CodeConflict)
	}
	member, err := lockedActiveMember(tx, hqID, memberID)
	if err != nil {
		return nil, nil, err
	}
	if member.PointsFrozen {
		return nil, nil, auth.NewError(auth.CodeConflict)
	}
	if member.PointsBalance > math.MaxInt32-points {
		return nil, nil, auth.NewError(auth.CodeValidationFailed)
	}
	return budget, member, nil
}

func (s *Service) persistGrant(tx *gorm.DB, principal *auth.WorkspacePrincipal, budget *gen.CustomerDailyPointGrantBudget, member *gen.CustomerMember, entry *gen.CustomerPointEntry) error {
	if err := tx.Create(entry).Error; err != nil {
		return err
	}
	if err := tx.Model(budget).Update("used_points", budget.UsedPoints+entry.Delta).Error; err != nil {
		return err
	}
	if err := tx.Model(member).Update("points_balance", member.PointsBalance+entry.Delta).Error; err != nil {
		return err
	}
	return s.auditPoint(tx, principal, entry)
}

func lockedGrantPolicy(tx *gorm.DB, hqID string, points int64) (*gen.CustomerBenefitPolicy, error) {
	var policy gen.CustomerBenefitPolicy
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("organization_id = ?", hqID).First(&policy).Error
	if err == gorm.ErrRecordNotFound {
		return nil, auth.NewError(auth.CodeConflict)
	}
	if err != nil {
		return nil, err
	}
	if policy.ManualGrantMaxSingle < points || policy.ManualGrantMaxDaily < points {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return &policy, nil
}

func (s *Service) lockDailyBudget(tx *gorm.DB, hqID string) (*gen.CustomerDailyPointGrantBudget, error) {
	now := s.now()
	date := now.In(shanghaiLocation).Format("2006-01-02")
	budget := &gen.CustomerDailyPointGrantBudget{
		ID: uuid.Must(uuid.NewV4()).String(), OrganizationID: hqID, BusinessDate: date, CreatedAt: now.UnixMilli(),
	}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(budget).Error; err != nil {
		return nil, err
	}
	budget.ID = ""
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("organization_id = ? AND business_date = ?", hqID, date).First(budget).Error
	return budget, err
}

func lockedActiveMember(tx *gorm.DB, hqID, memberID string) (*gen.CustomerMember, error) {
	var member gen.CustomerMember
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND organization_id = ?", memberID, hqID).First(&member).Error
	if err == gorm.ErrRecordNotFound {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err != nil {
		return nil, err
	}
	if member.Status != gen.CustomerMemberStatusActive {
		return nil, auth.NewError(auth.CodeConflict)
	}
	return &member, nil
}

func findPointByRequest(tx *gorm.DB, hqID string, kind gen.CustomerPointOperationKind, key string) (*gen.CustomerPointEntry, error) {
	var entry gen.CustomerPointEntry
	err := tx.Where("source_organization_id = ? AND operation_kind = ? AND request_key = ?", hqID, kind, key).First(&entry).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *Service) auditPoint(tx *gorm.DB, principal *auth.WorkspacePrincipal, entry *gen.CustomerPointEntry) error {
	return s.audit.Write(tx, audit.Entry{
		ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
		OrganizationID: principal.OrganizationID, Action: "hqCustomerPoints:" + strings.ToLower(string(entry.OperationKind)),
		ResourceType: "customerPointEntry", ResourceID: entry.ID, ResultCode: "SUCCESS",
		Metadata: audit.MetadataForPrincipal(principal, audit.Metadata{ReasonCode: string(entry.ReasonCode)}),
	})
}
