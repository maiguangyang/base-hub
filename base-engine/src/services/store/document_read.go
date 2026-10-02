package store

import (
	"context"
	"os"
	"path/filepath"

	"base-engine/auth"
	"base-engine/gen"
)

func (s *Service) OpenDocument(ctx context.Context, p *auth.WorkspacePrincipal, storeID, filename string) (*os.File, string, error) {
	if !validDocumentStoreID(storeID) || !validDocumentFilename(filename) {
		return nil, "", os.ErrNotExist
	}
	item, err := s.authorizedDocumentStore(s.db.WithContext(ctx), p, storeID, false)
	if err != nil {
		return nil, "", os.ErrNotExist
	}
	if !storeHasDocument(item, storeID, filename) {
		return nil, "", os.ErrNotExist
	}
	file, err := os.Open(s.storeDocumentPath(storeID, filename))
	if err != nil {
		return nil, "", err
	}
	return file, documentContentType(filename), nil
}

func storeHasDocument(item *gen.Store, storeID, filename string) bool {
	url := "/uploads/stores/" + storeID + "/" + filename
	return item.BusinessLicenseImageURL != nil && *item.BusinessLicenseImageURL == url ||
		item.OtherDocumentImageURL != nil && *item.OtherDocumentImageURL == url
}

func documentContentType(filename string) string {
	switch filepath.Ext(filename) {
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
