package productcatalog

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/uuid"
	"base-engine/auth"
	"base-engine/gen"
	"base-engine/src/services/audit"
	"base-engine/src/services/authorization"
	"gorm.io/gorm"
)

func (s *Service) SetProductMainImage(ctx context.Context, principal *auth.WorkspacePrincipal, productID, attachmentID string) (*gen.Product, error) {
	if !safeProductPathID(productID) {
		return nil, auth.NewError(auth.CodeValidationFailed)
	}
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	s.imageMu.Lock()
	defer s.imageMu.Unlock()
	pending, err := s.pendingForBinding(principal, hqID, attachmentID)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(s.imageRoot, "products", productID)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	filename := uuid.Must(uuid.NewV4()).String() + "." + pending.extension
	newPath := filepath.Join(dir, filename)
	pendingPath := s.pendingPath(attachmentID, pending.extension)
	item, oldURL, err := s.bindImageTx(ctx, principal, hqID, productID, pendingPath, newPath, filename)
	if err != nil {
		if _, statErr := os.Stat(newPath); statErr == nil {
			_ = os.Rename(newPath, pendingPath)
		}
		return nil, err
	}
	delete(s.pendingImages, attachmentID)
	if oldURL != "" {
		s.cleanupManagedImage(oldURL, productID)
	}
	return item, nil
}

func (s *Service) pendingForBinding(principal *auth.WorkspacePrincipal, hqID, attachmentID string) (pendingImage, error) {
	pending, ok := s.pendingImages[attachmentID]
	if !ok || time.Since(pending.createdAt) > time.Hour {
		return pendingImage{}, auth.NewError(auth.CodeValidationFailed)
	}
	if pending.accountID != principal.AccountID || pending.sessionID != principal.SessionID || pending.organizationID != hqID {
		return pendingImage{}, auth.NewError(auth.CodePermissionDenied)
	}
	return pending, nil
}

func (s *Service) bindImageTx(ctx context.Context, principal *auth.WorkspacePrincipal, hqID, productID, pendingPath, newPath, filename string) (*gen.Product, string, error) {
	var item gen.Product
	var oldURL string
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organization_id = ?", productID, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if item.ImageURL != nil {
			oldURL = *item.ImageURL
		}
		if err := os.Rename(pendingPath, newPath); err != nil {
			return err
		}
		newURL := "/uploads/products/" + productID + "/" + filename
		if err := tx.Model(&item).Update("image_url", newURL).Error; err != nil {
			return err
		}
		item.ImageURL = &newURL
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "product", ResourceID: productID, ResultCode: "SUCCESS"})
	})
	return &item, oldURL, err
}

func (s *Service) RemoveProductMainImage(ctx context.Context, principal *auth.WorkspacePrincipal, productID string) (*gen.Product, error) {
	hqID, err := headquartersID(principal, "hqProductCatalog:manage", authorization.AccessUpdate)
	if err != nil {
		return nil, err
	}
	var item gen.Product
	var oldURL string
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND organization_id = ?", productID, hqID).First(&item).Error; err != nil {
			return auth.NewError(auth.CodePermissionDenied)
		}
		if item.ImageURL != nil {
			oldURL = *item.ImageURL
		}
		if err := tx.Model(&item).Update("image_url", nil).Error; err != nil {
			return err
		}
		return s.audit.Write(tx, audit.Entry{ActorAccountID: principal.AccountID, SessionID: &principal.SessionID,
			OrganizationID: &hqID, Action: "hqProductCatalog:manage", ResourceType: "product", ResourceID: productID, ResultCode: "SUCCESS"})
	})
	if err != nil {
		return nil, err
	}
	if oldURL != "" {
		s.cleanupManagedImage(oldURL, productID)
	}
	item.ImageURL = nil
	return &item, nil
}

func (s *Service) cleanupManagedImage(imageURL, productID string) {
	path, ok := managedImagePath(s.imageRoot, productID, imageURL)
	if !ok {
		return
	}
	var count int64
	if err := s.db.Model(&gen.Product{}).Where("image_url = ?", imageURL).Count(&count).Error; err != nil || count != 0 {
		return
	}
	_ = os.Remove(path)
	_ = os.Remove(filepath.Dir(path))
}
