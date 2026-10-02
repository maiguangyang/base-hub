package storemerchandising

import (
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

type Service struct {
	db    *gorm.DB
	audit *audit.Service
}

func NewService(db *gorm.DB, auditService *audit.Service) *Service {
	return &Service{db: db, audit: auditService}
}

func storeScope(tx *gorm.DB, principal *auth.WorkspacePrincipal, storeID, action string, mode authorization.AccessMode) (*gen.Store, error) {
	if principal == nil {
		return nil, auth.NewError(auth.CodeAuthRequired)
	}
	if action == "" {
		return nil, auth.NewError(auth.CodePermissionDenied)
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

func scopedReadAction(principal *auth.WorkspacePrincipal, actions ...string) string {
	if principal == nil {
		return ""
	}
	for _, action := range actions {
		if principal.Has(action) {
			return action
		}
	}
	return ""
}
