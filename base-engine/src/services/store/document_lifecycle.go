package store

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func documentColumn(kind gen.StoreDocumentKind) (string, error) {
	switch kind {
	case gen.StoreDocumentKindBusinessLicense:
		return "business_license_image_url", nil
	case gen.StoreDocumentKindOther:
		return "other_document_image_url", nil
	default:
		return "", auth.NewError(auth.CodeValidationFailed)
	}
}

func documentURL(item *gen.Store, kind gen.StoreDocumentKind) *string {
	if kind == gen.StoreDocumentKindBusinessLicense {
		return item.BusinessLicenseImageURL
	}
	return item.OtherDocumentImageURL
}

func setDocumentURL(item *gen.Store, kind gen.StoreDocumentKind, url *string) {
	if kind == gen.StoreDocumentKindBusinessLicense {
		item.BusinessLicenseImageURL = url
	} else {
		item.OtherDocumentImageURL = url
	}
}

func (s *Service) authorizedDocumentStore(tx *gorm.DB, p *auth.WorkspacePrincipal, id string, write bool) (*gen.Store, error) {
	item, err := loadDocumentStore(tx, id, write)
	if err != nil {
		return nil, err
	}
	if err := authorizeStoreDocument(p, item, write); err != nil {
		return nil, err
	}
	return item, nil
}

func loadDocumentStore(tx *gorm.DB, id string, write bool) (*gen.Store, error) {
	item := &gen.Store{}
	query := tx
	if write {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := query.First(item, "id = ?", id).Error; err != nil {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	if item.IsDelete != nil && *item.IsDelete != 1 {
		return nil, auth.NewError(auth.CodePermissionDenied)
	}
	return item, nil
}

func authorizeStoreDocument(p *auth.WorkspacePrincipal, item *gen.Store, write bool) error {
	if p == nil || p.OrganizationID == nil {
		return auth.NewError(auth.CodeAuthRequired)
	}
	var action string
	mode := authorization.AccessRead
	if write {
		mode = authorization.AccessUpdate
	}
	switch p.WorkspaceType {
	case auth.WorkspaceTypeHeadquarters:
		action = headquartersDocumentAction(p, item, write)
	case auth.WorkspaceTypeFranchise:
		var err error
		action, err = franchiseDocumentAction(p, item, write)
		if err != nil {
			return err
		}
	default:
		return auth.NewError(auth.CodeWorkspaceForbidden)
	}
	return authorizeDocumentAction(p, item, action, mode)
}

func authorizeDocumentAction(p *auth.WorkspacePrincipal, item *gen.Store, action string, mode authorization.AccessMode) error {
	if action == "" {
		return auth.NewError(auth.CodePermissionDenied)
	}
	intent := authorization.Intent{Action: action, Mode: mode, ResourceOrganizationID: &item.OrganizationID}
	if p.WorkspaceType == auth.WorkspaceTypeFranchise && !p.AllStores {
		intent.StoreID = &item.ID
	}
	return authorization.Authorize(p, intent)
}

func headquartersDocumentAction(p *auth.WorkspacePrincipal, item *gen.Store, write bool) string {
	if !write && p.Has("store:read_all") {
		return "store:read_all"
	}
	if *p.OrganizationID != item.OrganizationID {
		return ""
	}
	if write {
		return "hqStore:update"
	}
	return "hqStore:read"
}

func franchiseDocumentAction(p *auth.WorkspacePrincipal, item *gen.Store, write bool) (string, error) {
	if *p.OrganizationID != item.OrganizationID {
		return "", auth.NewError(auth.CodePermissionDenied)
	}
	if write && item.Lifecycle != gen.StoreLifecycleDraft && item.Lifecycle != gen.StoreLifecycleRejected {
		return "", auth.NewError(auth.CodeConflict)
	}
	if write {
		return "store:update", nil
	}
	return "store:read", nil
}

func (s *Service) SetDocument(ctx context.Context, p *auth.WorkspacePrincipal, storeID string, kind gen.StoreDocumentKind, attachmentID string) (*gen.Store, error) {
	column, err := documentColumn(kind)
	if err != nil {
		return nil, err
	}
	if !validDocumentStoreID(storeID) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	s.documentMu.Lock()
	defer s.documentMu.Unlock()
	pending, err := s.pendingDocumentFor(p, attachmentID)
	if err != nil {
		return nil, err
	}
	filename := uuid.Must(uuid.NewV4()).String() + "." + pending.extension
	newPath := s.storeDocumentPath(storeID, filename)
	if err := os.MkdirAll(filepath.Dir(newPath), 0700); err != nil {
		return nil, err
	}
	pendingPath := s.pendingDocumentPath(attachmentID, pending.extension)
	item, oldURL, err := s.bindDocument(ctx, p, storeID, kind, column, pending, pendingPath, newPath, filename)
	if err != nil {
		if _, statErr := os.Stat(newPath); statErr == nil {
			_ = os.Rename(newPath, pendingPath)
		}
		return nil, err
	}
	delete(s.pendingDocuments, attachmentID)
	s.cleanupDocument(storeID, oldURL)
	return item, nil
}

func (s *Service) pendingDocumentFor(p *auth.WorkspacePrincipal, attachmentID string) (pendingDocument, error) {
	pending, ok := s.pendingDocuments[attachmentID]
	if !ok || time.Since(pending.createdAt) > time.Hour {
		return pendingDocument{}, auth.NewError(auth.CodeValidationFailed)
	}
	if p == nil || p.OrganizationID == nil || pending.accountID != p.AccountID || pending.sessionID != p.SessionID || pending.organizationID != *p.OrganizationID {
		return pendingDocument{}, auth.NewError(auth.CodePermissionDenied)
	}
	return pending, nil
}

func (s *Service) bindDocument(ctx context.Context, p *auth.WorkspacePrincipal, storeID string, kind gen.StoreDocumentKind, column string, pending pendingDocument, pendingPath, newPath, filename string) (*gen.Store, string, error) {
	var item *gen.Store
	var oldURL string
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		item, err = s.authorizedDocumentBindStore(tx, p, storeID, pending)
		if err != nil {
			return err
		}
		if old := documentURL(item, kind); old != nil {
			oldURL = *old
		}
		if err := os.Rename(pendingPath, newPath); err != nil {
			return err
		}
		url := "/uploads/stores/" + storeID + "/" + filename
		if err := tx.Model(item).Update(column, url).Error; err != nil {
			return err
		}
		setDocumentURL(item, kind, &url)
		return s.writeAudit(tx, p, item, "store:document:set", "")
	})
	return item, oldURL, err
}

func (s *Service) RemoveDocument(ctx context.Context, p *auth.WorkspacePrincipal, storeID string, kind gen.StoreDocumentKind) (*gen.Store, error) {
	column, err := documentColumn(kind)
	if err != nil {
		return nil, err
	}
	s.documentMu.Lock()
	defer s.documentMu.Unlock()
	var item *gen.Store
	var oldURL string
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		item, err = s.authorizedDocumentStore(tx, p, storeID, true)
		if err != nil {
			return err
		}
		if old := documentURL(item, kind); old != nil {
			oldURL = *old
		}
		if err := tx.Model(item).Update(column, nil).Error; err != nil {
			return err
		}
		setDocumentURL(item, kind, nil)
		return s.writeAudit(tx, p, item, "store:document:remove", "")
	})
	if err != nil {
		return nil, err
	}
	s.cleanupDocument(storeID, oldURL)
	return item, nil
}

func (s *Service) cleanupDocument(storeID, url string) {
	if url == "" {
		return
	}
	filename := filepath.Base(url)
	if !validDocumentFilename(filename) || url != "/uploads/stores/"+storeID+"/"+filename {
		return
	}
	var count int64
	if err := s.db.Model(&gen.Store{}).Where("business_license_image_url = ? OR other_document_image_url = ?", url, url).Count(&count).Error; err != nil || count != 0 {
		return
	}
	_ = os.Remove(s.storeDocumentPath(storeID, filename))
}
