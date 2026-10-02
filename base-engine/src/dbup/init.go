/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"time"

	"base-engine/src/services/ai"
	"base-engine/src/services/authentication"
	"gorm.io/gorm"
)

// SecurityBootstrap 保存一次性安全初始化占位，不包含任何秘密。
type SecurityBootstrap struct {
	Key       string    `gorm:"type:varchar(64);primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
}

type MembershipRoleJoin struct {
	OperatorMembershipID string `gorm:"column:operator_membership_id;type:varchar(36);primaryKey"`
	OperatorRoleID       string `gorm:"column:operator_role_id;type:varchar(36);primaryKey"`
}
func (MembershipRoleJoin) TableName() string { return "operator_membership_roles" }

type MembershipStoreJoin struct {
	OperatorMembershipID string `gorm:"column:operator_membership_id;type:varchar(36);primaryKey"`
	StoreID              string `gorm:"column:store_id;type:varchar(36);primaryKey"`
}
func (MembershipStoreJoin) TableName() string { return "operator_membership_stores" }

type PermissionRoleJoin struct {
	PermissionID   string `gorm:"column:permission_id;type:varchar(36);primaryKey"`
	OperatorRoleID string `gorm:"column:operator_role_id;type:varchar(36);primaryKey"`
}
func (PermissionRoleJoin) TableName() string { return "permission_roles" }

// MigrateSecurityTables 迁移不属于 GraphQL 模型的安全表与关联表。
func MigrateSecurityTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&authentication.AccountCredential{}, 
		&SecurityBootstrap{}, 
		&ai.StoredModelConfig{},
		&MembershipRoleJoin{},
		&MembershipStoreJoin{},
		&PermissionRoleJoin{},
	)
}
