/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"
	"fmt"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/gorm"
)

func requestPrincipal(ctx context.Context) (*auth.WorkspacePrincipal, error) {
	principal, err := auth.RequirePrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if err := authorizeInitialAccountQueryArgs(ctx, principal); err != nil {
		return nil, err
	}
	return principal, nil
}

func authorizeRequest(ctx context.Context, intent Intent) (*auth.WorkspacePrincipal, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	if err := Authorize(principal, intent); err != nil {
		return nil, err
	}
	return principal, nil
}

func resolverDB(ctx context.Context, resolver *gen.GeneratedResolver) (database *gorm.DB) {
	database = resolver.DB.Query().WithContext(ctx)
	defer func() { _ = recover() }()
	if transaction := gen.GetTransaction(ctx); transaction != nil {
		database = transaction
	}
	return database
}

func cloneInput(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(input)+2)
	for key, value := range input {
		result[key] = value
	}
	return result
}

func requireOrganization(principal *auth.WorkspacePrincipal) (string, error) {
	if principal.OrganizationID == nil || *principal.OrganizationID == "" {
		return "", auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return *principal.OrganizationID, nil
}

func denyEntityWrite(ctx context.Context) error {
	if _, err := requestPrincipal(ctx); err != nil {
		return err
	}
	return auth.NewError(auth.CodePermissionDenied)
}

func ensureInputOrganization(input map[string]interface{}, organizationID string) (map[string]interface{}, error) {
	result := cloneInput(input)
	removeGeneratedControlFields(result)
	if supplied, ok := result["organizationId"]; ok && supplied != nil && fmt.Sprint(supplied) != organizationID {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	delete(result, "organization")
	result["organizationId"] = organizationID
	return result, nil
}

func removeGeneratedControlFields(input map[string]interface{}) {
	for _, key := range []string{"isDelete", "state", "weight"} {
		delete(input, key)
	}
}

func activeRecordCondition(database *gorm.DB) *gorm.DB {
	return database.Where("is_delete IS NULL OR is_delete = ?", 1)
}

func membershipAccount(ctx context.Context, resolver *gen.GeneratedResolver, membership *gen.OperatorMembership) (*gen.Account, error) {
	if _, _, err := loadAuthorizedMembership(ctx, resolver, membership.ID, "read", AccessRelation); err != nil {
		return nil, err
	}
	item := &gen.Account{}
	return item, activeRecordCondition(resolver.DB.Query().WithContext(ctx)).First(item, "id = ?", membership.AccountID).Error
}

func membershipOrganization(ctx context.Context, resolver *gen.GeneratedResolver, membership *gen.OperatorMembership) (*gen.Organization, error) {
	if _, _, err := loadAuthorizedMembership(ctx, resolver, membership.ID, "read", AccessRelation); err != nil {
		return nil, err
	}
	item := &gen.Organization{}
	return item, activeRecordCondition(resolver.DB.Query().WithContext(ctx)).First(item, "id = ?", membership.OrganizationID).Error
}

func membershipRoles(ctx context.Context, resolver *gen.GeneratedResolver, membership *gen.OperatorMembership) ([]*gen.OperatorRole, error) {
	loaded, _, err := loadAuthorizedMembership(ctx, resolver, membership.ID, "read", AccessRelation)
	if err != nil {
		return nil, err
	}
	var items []*gen.OperatorRole
	err = resolver.DB.Query().WithContext(ctx).Model(loaded).
		Where("organization_id = ? AND (operator_roles.is_delete IS NULL OR operator_roles.is_delete = ?)", loaded.OrganizationID, 1).
		Association("Roles").Find(&items)
	return items, err
}

func membershipStores(ctx context.Context, resolver *gen.GeneratedResolver, membership *gen.OperatorMembership) ([]*gen.Store, error) {
	loaded, principal, err := loadAuthorizedMembership(ctx, resolver, membership.ID, "read", AccessRelation)
	if err != nil {
		return nil, err
	}
	var items []*gen.Store
	database := resolver.DB.Query().WithContext(ctx).Model(loaded).
		Where("organization_id = ? AND (stores.is_delete IS NULL OR stores.is_delete = ?)", loaded.OrganizationID, 1).
		Where("lifecycle = ?", gen.StoreLifecycleActive)
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise && !principal.AllStores {
		database = database.Where("stores.id IN ?", mapKeys(principal.StoreIDs))
	}
	err = database.Association("Stores").Find(&items)
	return items, err
}
