/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func pendingMembershipForReissue(tx *gorm.DB, accountID, organizationID string, now time.Time) (*gen.OperatorMembership, *gen.MembershipInvitation, error) {
	membership := &gen.OperatorMembership{}
	err := tx.Where("account_id = ? AND organization_id = ?", accountID, organizationID).First(membership).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if membership.Status != gen.MembershipStatusInvited || !recordActive(membership.IsDelete) {
		return nil, nil, auth.NewError(auth.CodeConflict)
	}
	invitation := &gen.MembershipInvitation{}
	err = tx.Where("membership_id = ?", membership.ID).First(invitation).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return membership, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}
	if invitation.AcceptedAt != nil || recordActive(invitation.IsDelete) && invitationUsable(invitation, now) {
		return nil, nil, auth.NewError(auth.CodeConflict)
	}
	return membership, invitation, nil
}

func resetPendingMembership(tx *gorm.DB, membership *gen.OperatorMembership, mode gen.StoreAccessMode, now time.Time) error {
	updates := map[string]any{
		"status": gen.MembershipStatusInvited, "store_access_mode": mode,
		"invited_at": now, "accepted_at": nil,
	}
	if err := tx.Model(membership).Updates(updates).Error; err != nil {
		return err
	}
	membership.Status = gen.MembershipStatusInvited
	membership.StoreAccessMode = mode
	membership.InvitedAt = &now
	membership.AcceptedAt = nil
	return nil
}

func refreshMembershipInvitation(tx *gorm.DB, invitation *gen.MembershipInvitation, membershipID, inviterID string, now time.Time) error {
	if invitation == nil {
		return tx.Create(&gen.MembershipInvitation{
			ID: uuid.Must(uuid.NewV4()).String(), MembershipID: membershipID,
			InvitedByAccountID: inviterID, ExpiresAt: now.Add(72 * time.Hour),
		}).Error
	}
	active := int64(1)
	updates := map[string]any{
		"invited_by_account_id": inviterID, "expires_at": now.Add(72 * time.Hour),
		"accepted_at": nil, "revoked_at": nil, "is_delete": active,
		"deleted_at": nil, "deleted_by": nil,
	}
	return tx.Model(invitation).Updates(updates).Error
}
