package customer

import (
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func couponStoreScope(tx *gorm.DB, principal *auth.WorkspacePrincipal, storeID, action string, mode authorization.AccessMode) (*gen.Store, error) {
	if principal == nil {
		return nil, auth.NewError(auth.CodeAuthRequired)
	}
	if principal.WorkspaceType != auth.WorkspaceTypeFranchise || principal.OrganizationID == nil {
		return nil, auth.NewError(auth.CodeWorkspaceForbidden)
	}
	var store gen.Store
	if err := tx.Where("id = ? AND organization_id = ?", storeID, *principal.OrganizationID).First(&store).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	intent := authorization.Intent{Action: action, Mode: mode, ResourceOrganizationID: &store.OrganizationID}
	if !principal.AllStores {
		intent.StoreID = &store.ID
	}
	if err := authorization.Authorize(principal, intent); err != nil {
		return nil, err
	}
	if mode != authorization.AccessRead {
		if err := authorization.AuthorizeOperationalStore(&store); err != nil {
			return nil, err
		}
	}
	return &store, nil
}

func couponHeadquartersID(tx *gorm.DB) (string, error) {
	var org gen.Organization
	if err := tx.Where("type = ? AND status = ?", gen.OrganizationTypeHeadquarters, gen.OrganizationStatusActive).First(&org).Error; err != nil {
		return "", err
	}
	return org.ID, nil
}
