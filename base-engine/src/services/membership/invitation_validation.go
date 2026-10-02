/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package membership

import (
	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func validateInvitationActivation(tx *gorm.DB, membership *gen.OperatorMembership) error {
	if membership.Status != gen.MembershipStatusInvited || !recordActive(membership.IsDelete) {
		return auth.NewError(auth.CodePermissionDenied)
	}
	var organizations int64
	err := tx.Model(&gen.Organization{}).
		Where("id = ? AND status = ?", membership.OrganizationID, gen.OrganizationStatusActive).
		Where("is_delete IS NULL OR is_delete = ?", 1).Count(&organizations).Error
	if err != nil || organizations != 1 {
		return auth.NewError(auth.CodePermissionDenied)
	}
	if err := validateInvitationRoles(tx, membership); err != nil {
		return err
	}
	return validateInvitationStores(tx, membership)
}

func validateInvitationRoles(tx *gorm.DB, membership *gen.OperatorMembership) error {
	total, err := associationCount(tx, "operator_membership_roles", "operator_membership_id", membership.ID)
	if err != nil {
		return err
	}
	var valid int64
	err = tx.Table("operator_membership_roles mr").
		Joins("JOIN operator_roles r ON r.id = mr.operator_role_id").
		Where("mr.operator_membership_id = ? AND r.organization_id = ?", membership.ID, membership.OrganizationID).
		Where("r.is_delete IS NULL OR r.is_delete = ?", 1).Count(&valid).Error
	if err != nil || valid != total {
		return auth.NewError(auth.CodePermissionDenied)
	}
	return nil
}

func validateInvitationStores(tx *gorm.DB, membership *gen.OperatorMembership) error {
	total, err := associationCount(tx, "operator_membership_stores", "operator_membership_id", membership.ID)
	if err != nil {
		return err
	}
	if membership.StoreAccessMode == gen.StoreAccessModeAllStores {
		if total == 0 {
			return nil
		}
		return auth.NewError(auth.CodeValidationFailed)
	}
	var valid int64
	err = tx.Table("operator_membership_stores ms").
		Joins("JOIN stores s ON s.id = ms.store_id").
		Where("ms.operator_membership_id = ? AND s.organization_id = ?", membership.ID, membership.OrganizationID).
		Where("s.lifecycle = ? AND (s.is_delete IS NULL OR s.is_delete = ?)", gen.StoreLifecycleActive, 1).
		Count(&valid).Error
	if err != nil || total == 0 || valid != total {
		return auth.NewError(auth.CodeStoreNotActive)
	}
	return nil
}

func associationCount(tx *gorm.DB, table, key, id string) (int64, error) {
	var count int64
	err := tx.Table(table).Where(key+" = ?", id).Count(&count).Error
	return count, err
}

func recordActive(isDelete *int64) bool {
	return isDelete == nil || *isDelete == 1
}
