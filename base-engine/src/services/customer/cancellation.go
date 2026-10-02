package customer

import (
	"context"
	"strings"
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) SetMemberStatus(ctx context.Context, principal *auth.WorkspacePrincipal, memberID string, status gen.CustomerMemberStatus) (*MemberView, error) {
	hqID, err := headquartersID(principal, "hqCustomer:update", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	if !mutableCustomerMemberStatus(status) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var member gen.CustomerMember
	var pending int64
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockMember(tx, hqID, memberID, &member); err != nil {
			return err
		}
		if member.Status != gen.CustomerMemberStatusActive && member.Status != gen.CustomerMemberStatusSuspended {
			return auth.NewError(auth.CodeConflict)
		}
		previousStatus := member.Status
		var err error
		pending, err = countPendingCoupons(tx, member.ID)
		if err != nil {
			return err
		}
		if err := tx.Model(&member).Update("status", status).Error; err != nil {
			return err
		}
		member.Status = status
		if err := scheduleReactivatedMemberCouponCatchup(tx, &member, previousStatus, s.now().UnixMilli()); err != nil {
			return err
		}
		return s.auditMember(tx, principal, "hqCustomer:update", member.ID, "")
	})
	if err != nil {
		return nil, err
	}
	view := viewMember(&member)
	view.PendingCouponCount = pending
	return &view, nil
}

func mutableCustomerMemberStatus(status gen.CustomerMemberStatus) bool {
	return status == gen.CustomerMemberStatusActive || status == gen.CustomerMemberStatusSuspended
}

func (s *Service) RequestCancellation(ctx context.Context, principal *auth.WorkspacePrincipal, memberID, identityEvidence, basisCode string) (*MemberView, error) {
	hqID, err := headquartersID(principal, "hqCustomer:cancel", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(identityEvidence) == "" || strings.TrimSpace(basisCode) == "" {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var member gen.CustomerMember
	var pending int64
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockMember(tx, hqID, memberID, &member); err != nil {
			return err
		}
		if member.Status != gen.CustomerMemberStatusActive && member.Status != gen.CustomerMemberStatusSuspended {
			return auth.NewError(auth.CodeConflict)
		}
		var err error
		pending, err = countPendingCoupons(tx, member.ID)
		if err != nil {
			return err
		}
		now := time.Now()
		if err := tx.Model(&member).Updates(map[string]any{
			"status":                    gen.CustomerMemberStatusCancelPending,
			"cancellation_requested_at": now,
		}).Error; err != nil {
			return err
		}
		member.Status, member.CancellationRequestedAt = gen.CustomerMemberStatusCancelPending, &now
		return s.auditMember(tx, principal, "hqCustomer:cancel_request", member.ID, identityEvidence, basisCode)
	})
	if err != nil {
		return nil, err
	}
	view := viewMember(&member)
	view.PendingCouponCount = pending
	return &view, nil
}

func (s *Service) CompleteCancellation(ctx context.Context, principal *auth.WorkspacePrincipal, memberID, dispositionReference string) (*gen.CustomerMember, error) {
	hqID, err := headquartersID(principal, "hqCustomer:cancel", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(dispositionReference) == "" {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	var member gen.CustomerMember
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockMember(tx, hqID, memberID, &member); err != nil {
			return err
		}
		if member.Status != gen.CustomerMemberStatusCancelPending {
			return auth.NewError(auth.CodeConflict)
		}
		if err := ensureCancellationPointsSettled(tx, &member); err != nil {
			return err
		}
		pending, err := countPendingCoupons(tx, member.ID)
		if err != nil {
			return err
		}
		if pending != 0 {
			return auth.NewError(auth.CodeConflict)
		}
		now := time.Now()
		if err := tx.Model(&member).Updates(map[string]any{
			"status": gen.CustomerMemberStatusCancelled, "phone": nil, "cancelled_at": now,
		}).Error; err != nil {
			return err
		}
		member.Status, member.Phone, member.CancelledAt = gen.CustomerMemberStatusCancelled, nil, &now
		return s.auditMember(tx, principal, "hqCustomer:cancel_complete", member.ID, dispositionReference)
	})
	return &member, err
}

func ensureCancellationPointsSettled(tx *gorm.DB, member *gen.CustomerMember) error {
	if member.PointsBalance != 0 || member.PointsFrozen {
		return auth.NewError(auth.CodeConflict)
	}
	ledgerPoints, err := pointLedgerSum(tx, member.ID)
	if err != nil {
		return err
	}
	if ledgerPoints != 0 {
		return auth.NewError(auth.CodeConflict)
	}
	return nil
}

func countPendingCoupons(tx *gorm.DB, memberID string) (int64, error) {
	var pending int64
	err := tx.Model(&gen.CustomerCouponGrant{}).Where(
		"member_id = ? AND status = ?", memberID, gen.CustomerCouponGrantStatusPendingActivation,
	).Count(&pending).Error
	return pending, err
}

func lockMember(tx *gorm.DB, hqID, memberID string, target *gen.CustomerMember) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(
		"organization_id = ? AND (is_delete IS NULL OR is_delete = 1)", hqID,
	).First(target, "id = ?", memberID).Error
	if err == gorm.ErrRecordNotFound {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return err
}
