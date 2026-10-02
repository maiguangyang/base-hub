/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package authorization

import (
	"context"

	"base-engine/auth"
	"base-engine/gen"
)

func secureMembershipInput(ctx context.Context, resolver *gen.GeneratedResolver, principal *auth.WorkspacePrincipal, input map[string]interface{}, organizationID string) (map[string]interface{}, error) {
	if err := rejectNestedMembershipRelations(input); err != nil {
		return nil, err
	}
	if err := validateMembershipReferences(ctx, resolver, principal, input, organizationID); err != nil {
		return nil, err
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		if err := ValidateHQRoleDelegation(resolverDB(ctx, resolver), principal, organizationID, stringIDs(input["rolesIds"])); err != nil {
			return nil, err
		}
		input["storeAccessMode"] = gen.StoreAccessModeAllStores
		input["storesIds"] = []string{}
	}
	normalizeMembershipEnums(input)
	if err := validateMembershipModeScope(principal, input); err != nil {
		return nil, err
	}
	return input, nil
}

func rejectNestedMembershipRelations(input map[string]interface{}) error {
	for _, key := range []string{"account", "organization", "roles", "stores", "invitations", "invitationsIds"} {
		if input[key] != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
	}
	return nil
}

func enforceMembershipStoreMode(ctx context.Context, resolver *gen.GeneratedResolver, membership *gen.OperatorMembership, input map[string]interface{}) error {
	mode := membership.StoreAccessMode
	if supplied, ok := input["storeAccessMode"].(gen.StoreAccessMode); ok {
		mode = supplied
	}
	storeIDsValue, storesSupplied := input["storesIds"]
	stores := stringIDs(storeIDsValue)
	if mode == gen.StoreAccessModeAllStores {
		if storesSupplied && len(stores) > 0 {
			return auth.NewError(auth.CodeValidationFailed)
		}
		input["storesIds"] = []string{}
		return nil
	}
	if mode != gen.StoreAccessModeSelectedStores {
		return auth.NewError(auth.CodeValidationFailed)
	}
	if storesSupplied {
		if len(stores) == 0 {
			return auth.NewError(auth.CodeValidationFailed)
		}
		return nil
	}
	association := resolverDB(ctx, resolver).Model(membership).Association("Stores")
	count := association.Count()
	if association.Error != nil {
		return association.Error
	}
	if count == 0 {
		return auth.NewError(auth.CodeValidationFailed)
	}
	return nil
}

func validateMembershipReferences(ctx context.Context, resolver *gen.GeneratedResolver, principal *auth.WorkspacePrincipal, input map[string]interface{}, organizationID string) error {
	checks := []struct {
		key, table string
	}{
		{key: "rolesIds", table: "operator_roles"},
		{key: "storesIds", table: "stores"},
	}
	for _, check := range checks {
		ids := stringIDs(input[check.key])
		if check.table == "stores" {
			if err := validatePrincipalStoreScope(principal, ids); err != nil {
				return err
			}
		}
		if err := validateOrganizationIDs(ctx, resolver, check.table, ids, organizationID); err != nil {
			return err
		}
	}
	return nil
}

func validateMembershipModeScope(principal *auth.WorkspacePrincipal, input map[string]interface{}) error {
	if mode, ok := input["storeAccessMode"].(gen.StoreAccessMode); ok && mode == gen.StoreAccessModeAllStores && principal.WorkspaceType == auth.WorkspaceTypeFranchise && !principal.AllStores {
		return auth.NewError(auth.CodeStoreScopeDenied)
	}
	return nil
}

func validatePrincipalStoreScope(principal *auth.WorkspacePrincipal, ids []string) error {
	if principal == nil || principal.WorkspaceType != auth.WorkspaceTypeFranchise || principal.AllStores {
		return nil
	}
	for _, id := range ids {
		if !principal.HasStore(id) {
			return auth.NewError(auth.CodeStoreScopeDenied)
		}
	}
	return nil
}

func validateOrganizationIDs(ctx context.Context, resolver *gen.GeneratedResolver, table string, ids []string, organizationID string) error {
	if len(ids) == 0 {
		return nil
	}
	query := resolverDB(ctx, resolver).Table(table).
		Where("id IN ? AND organization_id = ?", ids, organizationID).
		Where("is_delete IS NULL OR is_delete = ?", 1)
	code := auth.CodePermissionDenied
	if table == "stores" {
		query = query.Where("lifecycle = ?", gen.StoreLifecycleActive)
		code = auth.CodeStoreNotActive
	}
	var count int64
	err := query.Count(&count).Error
	if err != nil || count != int64(len(ids)) {
		return auth.NewError(code)
	}
	return nil
}

func normalizeMembershipEnums(input map[string]interface{}) {
	if mode, ok := input["storeAccessMode"].(*gen.StoreAccessMode); ok && mode != nil {
		input["storeAccessMode"] = *mode
	}
}
