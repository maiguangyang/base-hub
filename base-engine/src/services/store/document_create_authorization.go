package store

import (
	"time"

	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) authorizedDocumentBindStore(tx *gorm.DB, p *auth.WorkspacePrincipal, id string, pending pendingDocument) (*gen.Store, error) {
	item, err := loadDocumentStore(tx, id, true)
	if err != nil {
		return nil, err
	}
	err = authorizeStoreDocument(p, item, true)
	if err == nil {
		return item, nil
	}
	if auth.ErrorCode(err) != auth.CodePermissionDenied || !createdWithPendingDocument(p, item, pending) {
		return nil, err
	}
	return item, authorizeDocumentCreate(p, item)
}

func createdWithPendingDocument(p *auth.WorkspacePrincipal, item *gen.Store, pending pendingDocument) bool {
	if p == nil || p.OrganizationID == nil || item.CreatedBy == nil || *item.CreatedBy != p.AccountID || item.OrganizationID != *p.OrganizationID {
		return false
	}
	if p.WorkspaceType == auth.WorkspaceTypeFranchise && item.Lifecycle != gen.StoreLifecycleDraft {
		return false
	}
	start := pending.createdAt.UnixMilli()
	return item.CreatedAt >= start && item.CreatedAt <= start+int64(time.Hour/time.Millisecond)
}

func authorizeDocumentCreate(p *auth.WorkspacePrincipal, item *gen.Store) error {
	action := "store:create"
	if p.WorkspaceType == auth.WorkspaceTypeHeadquarters {
		action = "hqStore:create"
	}
	return authorization.Authorize(p, authorization.Intent{Action: action, Mode: authorization.AccessCreate, ResourceOrganizationID: &item.OrganizationID})
}
