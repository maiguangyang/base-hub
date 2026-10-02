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
	"gorm.io/gorm/clause"
)

func registerStoreHandlers(handlers *gen.ResolutionHandlers) {
	handlers.CreateStore = createStore
	handlers.UpdateStore = updateStore
	handlers.DeleteStores = deleteStores
	handlers.RecoveryStores = recoverStores
	handlers.QueryStore = queryStore
	handlers.QueryStores = queryStores
	handlers.StoreOrganization = storeOrganization
	handlers.StoreMembers = storeMembers
	handlers.StoreReviewedByAccount = storeReviewedByAccount
	handlers.StoreAuditLogs = storeAuditLogs
}

func storeAction(principal *auth.WorkspacePrincipal, verb string) string {
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		if verb == "read" {
			return "store:read_all"
		}
		return "hqStore:" + verb
	}
	return "store:" + verb
}

func createStore(ctx context.Context, resolver *gen.GeneratedResolver, input map[string]interface{}) (*gen.Store, error) {
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, err
	}
	organizationID, err := requireOrganization(principal)
	if err != nil {
		return nil, err
	}
	if err := Authorize(principal, Intent{Action: storeAction(principal, "create"), Mode: AccessCreate, ResourceOrganizationID: &organizationID}); err != nil {
		return nil, err
	}
	secured, err := ensureInputOrganization(input, organizationID)
	if err != nil {
		return nil, err
	}
	if err := rejectStoreManagedRelations(secured); err != nil {
		return nil, err
	}
	secured["lifecycle"] = gen.StoreLifecycleDraft
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		secured["lifecycle"] = gen.StoreLifecycleActive
	}
	removeStoreGovernanceFields(secured)
	secured["lifecycle"] = initialStoreLifecycle(principal)
	normalizeStoreBusinessStatus(secured)
	created, err := gen.CreateStoreHandler(context.WithValue(ctx, gen.KeyPrincipalID, &principal.AccountID), resolver, secured)
	if err != nil {
		return nil, err
	}
	return created, auditStoreWrite(ctx, resolver, principal, created, "create")
}

func updateStore(ctx context.Context, resolver *gen.GeneratedResolver, id string, input map[string]interface{}) (*gen.Store, error) {
	store, principal, err := loadAuthorizedStore(ctx, resolver, id, "update", AccessUpdate)
	if err != nil {
		return nil, err
	}
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise && store.Lifecycle != gen.StoreLifecycleDraft && store.Lifecycle != gen.StoreLifecycleRejected {
		return nil, auth.NewError(auth.CodeConflict)
	}
	secured, err := ensureInputOrganization(input, store.OrganizationID)
	if err != nil {
		return nil, err
	}
	if err := rejectStoreManagedRelations(secured); err != nil {
		return nil, err
	}
	removeStoreGovernanceFields(secured)
	normalizeStoreBusinessStatus(secured)
	updated, err := gen.UpdateStoreHandler(ctx, resolver, id, secured)
	if err != nil {
		return nil, err
	}
	return updated, auditStoreWrite(ctx, resolver, principal, updated, "update")
}

func loadAuthorizedStore(ctx context.Context, resolver *gen.GeneratedResolver, id, verb string, mode AccessMode) (*gen.Store, *auth.WorkspacePrincipal, error) {
	store := &gen.Store{}
	writeMode := mode != AccessRead && mode != AccessRelation
	database := resolverDB(ctx, resolver)
	if writeMode {
		database = database.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := database.First(store, "id = ?", id).Error; err != nil {
		return nil, nil, auth.NewError(auth.CodePermissionDenied)
	}
	if err := authorizeDeletionState(store.IsDelete, mode); err != nil {
		return nil, nil, err
	}
	principal, err := requestPrincipal(ctx)
	if err != nil {
		return nil, nil, err
	}
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters && writeMode && (principal.OrganizationID == nil || *principal.OrganizationID != store.OrganizationID) {
		return nil, nil, auth.NewError(auth.CodePermissionDenied)
	}
	action := storeResourceAction(principal, store.OrganizationID, verb)
	err = Authorize(principal, Intent{Action: action, Mode: mode, ResourceOrganizationID: &store.OrganizationID, StoreID: storeScopeID(principal, store.ID)})
	return store, principal, err
}

func storeResourceAction(principal *auth.WorkspacePrincipal, organizationID, verb string) string {
	if verb == "read" {
		return storeReadAction(principal, organizationID)
	}
	return storeAction(principal, verb)
}

func storeReadAction(principal *auth.WorkspacePrincipal, organizationID string) string {
	if principal.WorkspaceType != auth.WorkspaceTypeHeadquarters {
		return "store:read"
	}
	if principal.Has("store:read_all") {
		return "store:read_all"
	}
	if principal.OrganizationID != nil && *principal.OrganizationID == organizationID {
		return "hqStore:read"
	}
	return "store:read_all"
}

func storeScopeID(principal *auth.WorkspacePrincipal, storeID string) *string {
	if principal.WorkspaceType == auth.WorkspaceTypeFranchise && !principal.AllStores {
		return &storeID
	}
	return nil
}

// gqlgen returns optional enums as pointers; ApplyChanges expects their wire strings.
func normalizeStoreBusinessStatus(input map[string]interface{}) {
	if status, ok := input["businessStatus"].(*gen.StoreBusinessStatus); ok && status != nil {
		input["businessStatus"] = status.String()
	}
}

func initialStoreLifecycle(principal *auth.WorkspacePrincipal) gen.StoreLifecycle {
	if principal.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		return gen.StoreLifecycleActive
	}
	return gen.StoreLifecycleDraft
}

func removeStoreGovernanceFields(input map[string]interface{}) {
	for _, key := range []string{"lifecycle", "submittedAt", "reviewedAt", "reviewedByAccountId", "rejectionReason"} {
		delete(input, key)
	}
}

func rejectStoreManagedRelations(input map[string]interface{}) error {
	if err := rejectStoreDocumentFields(input); err != nil {
		return err
	}
	for _, key := range []string{"members", "membersIds", "reviewedByAccount", "reviewedByAccountId", "auditLogs", "auditLogsIds", "paymentConfigs", "paymentConfigsIds", "stocktakes", "stocktakesIds"} {
		if input[key] != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
	}
	return nil
}

func storeOrganization(ctx context.Context, resolver *gen.GeneratedResolver, store *gen.Store) (*gen.Organization, error) {
	if _, _, err := loadAuthorizedStore(ctx, resolver, store.ID, "read", AccessRelation); err != nil {
		return nil, err
	}
	item := &gen.Organization{}
	return item, activeRecordCondition(resolver.DB.Query().WithContext(ctx)).First(item, "id = ?", store.OrganizationID).Error
}

func storeMembers(ctx context.Context, resolver *gen.GeneratedResolver, store *gen.Store) ([]*gen.OperatorMembership, error) {
	loaded, _, err := loadAuthorizedStore(ctx, resolver, store.ID, "read", AccessRelation)
	if err != nil {
		return nil, err
	}
	var items []*gen.OperatorMembership
	err = resolver.DB.Query().WithContext(ctx).Model(loaded).
		Where("organization_id = ? AND (operator_memberships.is_delete IS NULL OR operator_memberships.is_delete = ?)", loaded.OrganizationID, 1).
		Association("Members").Find(&items)
	return items, err
}

func mapKeys(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	return result
}
