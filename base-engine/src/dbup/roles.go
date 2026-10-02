/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"errors"

	"github.com/gofrs/uuid"
	"base-engine/gen"
	"gorm.io/gorm"
)

const (
	HeadquartersCode     = "HQ"
	HeadquartersName     = "organization.headquarters"
	HeadquartersRoleName = "role.hqSuperAdministrator"
)

var errHeadquartersNotUnique = errors.New("HQ_ORGANIZATION_NOT_UNIQUE")

// InitRoles 确保系统仅有一个总部组织及其动态超级管理员角色。
func InitRoles(db *gorm.DB) error {
	var organizations []gen.Organization
	if err := db.Where("type = ?", gen.OrganizationTypeHeadquarters).Find(&organizations).Error; err != nil {
		return err
	}
	if len(organizations) > 1 {
		return errHeadquartersNotUnique
	}
	organization, err := ensureHeadquartersOrganization(db, organizations)
	if err != nil {
		return err
	}
	return ensureHeadquartersRole(db, organization.ID)
}

// InitRolePermissions 保留静态自定义角色种子入口；系统角色按 scope 动态解析。
func InitRolePermissions(_ *gorm.DB) error {
	return nil
}

// ResolveDynamicPermissions 按系统角色种类解析最新权限，无需物化关联行。
func ResolveDynamicPermissions(kind gen.RoleKind, seeds []PermissionSeed) []PermissionSeed {
	wanted := dynamicRoleScope(kind)
	result := make([]PermissionSeed, 0, len(seeds))
	for _, seed := range seeds {
		if seed.Scope == wanted {
			result = append(result, seed)
		}
	}
	return result
}

func dynamicRoleScope(kind gen.RoleKind) gen.PermissionScope {
	if kind == gen.RoleKindHqSuperAdmin {
		return gen.PermissionScopeSystem
	}
	if kind == gen.RoleKindFranchiseOwner {
		return gen.PermissionScopeTenant
	}
	return ""
}

func ensureHeadquartersOrganization(db *gorm.DB, organizations []gen.Organization) (gen.Organization, error) {
	if len(organizations) == 1 {
		return organizations[0], nil
	}
	organization := gen.Organization{
		ID: uuid.Must(uuid.NewV4()).String(), Code: HeadquartersCode,
		Name: HeadquartersName, Type: gen.OrganizationTypeHeadquarters,
		Status: gen.OrganizationStatusActive,
	}
	return organization, db.Create(&organization).Error
}

func ensureHeadquartersRole(db *gorm.DB, organizationID string) error {
	role := gen.OperatorRole{}
	err := db.Where("organization_id = ? AND kind = ?", organizationID, gen.RoleKindHqSuperAdmin).First(&role).Error
	if err == nil {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	role = gen.OperatorRole{
		ID: uuid.Must(uuid.NewV4()).String(), Name: HeadquartersRoleName,
		Kind: gen.RoleKindHqSuperAdmin, OrganizationID: organizationID,
	}
	return db.Create(&role).Error
}
