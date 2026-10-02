/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"strings"

	"github.com/gofrs/uuid"
	"base-engine/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PermissionSeed 是后端权威权限目录项。
type PermissionSeed struct {
	Name   string
	Action string
	Module string
	Scope  gen.PermissionScope
}

var systemCRUDResources = []string{
	"account", "organization", "permission", "session", "auditLog", "franchiseOpeningRecord",
	"globalPaymentConfig", "franchisePaymentConfig", "storePaymentConfig",
}
var tenantCRUDResources = []string{
	"operatorMembership", "operatorRole", "store", "membershipInvitation",
}

var systemActions = []string{
	"franchise:provision", "organization:suspend", "organization:restore",
	"hqRole:read", "hqRole:create", "hqRole:update", "hqRole:delete",
	"hqMembership:read", "hqMembership:create", "hqMembership:update", "hqMembership:delete",
	"hqStore:read", "hqStore:create", "hqStore:update", "hqStore:delete",
	"store:read_all", "store:approve", "store:reject", "report:export",
	"aiModelConfig:read", "aiModelConfig:manage",
	"paymentConfig:read", "paymentConfig:manage",
}

var tenantActions = []string{"tenantAudit:read", "store:submit"}

// PermissionSeeds 返回完整且确定有序的权限目录。
func PermissionSeeds() []PermissionSeed {
	seeds := make([]PermissionSeed, 0, 58)
	seeds = appendCRUDSeeds(seeds, systemCRUDResources, gen.PermissionScopeSystem)
	seeds = appendCRUDSeeds(seeds, tenantCRUDResources, gen.PermissionScopeTenant)
	seeds = appendActionSeeds(seeds, systemActions, gen.PermissionScopeSystem)
	return appendActionSeeds(seeds, tenantActions, gen.PermissionScopeTenant)
}

// InitPermissions 以 action 为自然键幂等写入权威目录。
func InitPermissions(db *gorm.DB) error {
	for _, seed := range PermissionSeeds() {
		permission := seed.toModel()
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "action"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "module", "scope"}),
		}).Create(&permission).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// NormalizePermissionScopes 修复既有权限目录中的作用域漂移。
func NormalizePermissionScopes(db *gorm.DB) error {
	for _, seed := range PermissionSeeds() {
		err := db.Model(&gen.Permission{}).Where("action = ?", seed.Action).Updates(map[string]any{
			"name": seed.Name, "module": seed.Module, "scope": seed.Scope,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func appendCRUDSeeds(seeds []PermissionSeed, resources []string, scope gen.PermissionScope) []PermissionSeed {
	for _, resource := range resources {
		for _, verb := range []string{"read", "create", "update", "delete"} {
			seeds = append(seeds, newPermissionSeed(resource+":"+verb, scope))
		}
	}
	return seeds
}

func appendActionSeeds(seeds []PermissionSeed, actions []string, scope gen.PermissionScope) []PermissionSeed {
	for _, action := range actions {
		seeds = append(seeds, newPermissionSeed(action, scope))
	}
	return seeds
}

func newPermissionSeed(action string, scope gen.PermissionScope) PermissionSeed {
	parts := strings.SplitN(action, ":", 2)
	return PermissionSeed{
		Name:   "permission." + parts[0] + "." + parts[1],
		Action: action, Module: parts[0], Scope: scope,
	}
}

func (seed PermissionSeed) toModel() gen.Permission {
	return gen.Permission{
		ID: uuid.Must(uuid.NewV4()).String(), Name: seed.Name,
		Action: seed.Action, Module: seed.Module, Scope: seed.Scope,
	}
}
