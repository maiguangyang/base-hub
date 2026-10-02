/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package dbup

import (
	"strings"
	"testing"

	"base-engine/gen"
)

// TestPermissionSeedContract 锁定 CRUD、扩展动作、作用域与动态系统角色语义。
func TestPermissionSeedContract(t *testing.T) {
	seeds := PermissionSeeds()
	byAction := make(map[string]PermissionSeed, len(seeds))
	for _, seed := range seeds {
		if _, exists := byAction[seed.Action]; exists {
			t.Fatalf("duplicate action: %s", seed.Action)
		}
		byAction[seed.Action] = seed
		if seed.Name != "permission."+strings.Replace(seed.Action, ":", ".", 1) {
			t.Fatalf("invalid localization key: %s", seed.Name)
		}
	}
	assertEntityCRUDScopes(t, byAction)
	assertAdditionalActions(t, byAction)
	assertDynamicRoleScopes(t, seeds)
}

func TestPaymentConfigPermissionSeeds(t *testing.T) {
	byAction := make(map[string]PermissionSeed)
	for _, seed := range PermissionSeeds() {
		byAction[seed.Action] = seed
	}
	for _, action := range []string{"paymentConfig:read", "paymentConfig:manage", "globalPaymentConfig:read", "globalPaymentConfig:create", "franchisePaymentConfig:read", "franchisePaymentConfig:create", "storePaymentConfig:read", "storePaymentConfig:create"} {
		seed, ok := byAction[action]
		if !ok || seed.Scope != gen.PermissionScopeSystem {
			t.Fatalf("missing system payment permission %s", action)
		}
	}
}

func TestCouponDistributionJobPermissionSeeds(t *testing.T) {
	byAction := make(map[string]PermissionSeed)
	for _, seed := range PermissionSeeds() {
		byAction[seed.Action] = seed
	}
	for _, verb := range []string{"read", "create", "update", "delete"} {
		action := "customerCouponDistributionJob:" + verb
		seed, ok := byAction[action]
		if !ok || seed.Scope != gen.PermissionScopeSystem {
			t.Fatalf("missing system job permission %s", action)
		}
	}
}

// TestPermissionSeedsAreIdempotentAndNormalized 验证重复启动不增行并修复作用域漂移。
func TestPermissionSeedsAreIdempotentAndNormalized(t *testing.T) {
	db := openTestDB(t)
	if err := migrateGeneratedForSQLite(db, &gen.Permission{}); err != nil {
		t.Fatal(err)
	}
	if err := InitPermissions(db); err != nil {
		t.Fatal(err)
	}
	if err := InitPermissions(db); err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := db.Model(&gen.Permission{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != int64(len(PermissionSeeds())) {
		t.Fatalf("permission count = %d", count)
	}
	if err := db.Model(&gen.Permission{}).Where("action = ?", "store:read").Update("scope", gen.PermissionScopeSystem).Error; err != nil {
		t.Fatal(err)
	}
	if err := NormalizePermissionScopes(db); err != nil {
		t.Fatal(err)
	}
	var permission gen.Permission
	if err := db.Where("action = ?", "store:read").First(&permission).Error; err != nil {
		t.Fatal(err)
	}
	if permission.Scope != gen.PermissionScopeTenant {
		t.Fatalf("scope was not normalized: %s", permission.Scope)
	}
}

func assertEntityCRUDScopes(t *testing.T, seeds map[string]PermissionSeed) {
	t.Helper()
	system := map[string]bool{"account": true, "organization": true, "permission": true, "session": true, "auditLog": true, "franchiseOpeningRecord": true,
		"customerMember": true, "customerBenefitPolicy": true, "customerDailyPointGrantBudget": true,
		"customerPointEntry": true, "customerCouponTemplate": true, "customerCouponGrant": true, "customerCouponDistributionJob": true}
	resources := []string{"account", "organization", "operatorMembership", "permission", "operatorRole", "store", "session", "membershipInvitation", "auditLog", "franchiseOpeningRecord",
		"customerMember", "customerBenefitPolicy", "customerDailyPointGrantBudget", "customerPointEntry", "customerCouponTemplate", "customerCouponGrant", "customerCouponDistributionJob"}
	for _, resource := range resources {
		expected := gen.PermissionScopeTenant
		if system[resource] {
			expected = gen.PermissionScopeSystem
		}
		for _, verb := range []string{"read", "create", "update", "delete"} {
			seed, ok := seeds[resource+":"+verb]
			if !ok || seed.Scope != expected {
				t.Fatalf("invalid CRUD seed %s:%s: %#v", resource, verb, seed)
			}
		}
	}
}

func assertAdditionalActions(t *testing.T, seeds map[string]PermissionSeed) {
	t.Helper()
	system := []string{"franchise:provision", "organization:suspend", "organization:restore", "hqRole:read", "hqRole:create", "hqRole:update", "hqRole:delete", "hqMembership:read", "hqMembership:create", "hqMembership:update", "hqMembership:delete", "hqStore:read", "hqStore:create", "hqStore:update", "hqStore:delete", "store:read_all", "store:approve", "store:reject", "customer:read_sensitive", "report:export", "aiModelConfig:read", "aiModelConfig:manage",
		"hqCustomer:read", "hqCustomer:create", "hqCustomer:update", "hqCustomer:cancel",
		"hqCustomerPolicy:read", "hqCustomerPolicy:manage",
		"hqCustomerPoints:read", "hqCustomerPoints:grant", "hqCustomerPoints:reverse", "hqCustomerPoints:correct",
		"hqCustomerCoupon:read", "hqCustomerCoupon:manage", "hqCustomerCoupon:grant", "hqCustomerCoupon:revoke"}
	for _, action := range system {
		if seeds[action].Scope != gen.PermissionScopeSystem {
			t.Fatalf("%s must be SYSTEM", action)
		}
	}
	for _, action := range []string{"tenantAudit:read", "store:submit"} {
		if seeds[action].Scope != gen.PermissionScopeTenant {
			t.Fatalf("%s must be TENANT", action)
		}
	}
}

func assertDynamicRoleScopes(t *testing.T, seeds []PermissionSeed) {
	t.Helper()
	owner := ResolveDynamicPermissions(gen.RoleKindFranchiseOwner, seeds)
	superAdmin := ResolveDynamicPermissions(gen.RoleKindHqSuperAdmin, seeds)
	for _, seed := range seeds {
		if (seed.Scope == gen.PermissionScopeTenant) != containsAction(owner, seed.Action) {
			t.Fatalf("owner dynamic scope mismatch: %s", seed.Action)
		}
		if (seed.Scope == gen.PermissionScopeSystem) != containsAction(superAdmin, seed.Action) {
			t.Fatalf("super admin dynamic scope mismatch: %s", seed.Action)
		}
	}
}

func containsAction(seeds []PermissionSeed, action string) bool {
	for _, seed := range seeds {
		if seed.Action == action {
			return true
		}
	}
	return false
}
