package store

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"base-engine/gen"
)

func (s *Service) SweepDocumentFiles(ctx context.Context) error {
	s.documentMu.Lock()
	defer s.documentMu.Unlock()
	now := time.Now()
	s.sweepPendingDocuments(now)
	if err := s.sweepPendingDocumentDirectory(now); err != nil {
		return err
	}
	return s.sweepBoundDocumentDirectory(ctx)
}

func (s *Service) sweepPendingDocumentDirectory(now time.Time) error {
	dir := filepath.Join(s.documentRoot, "store-documents", "pending")
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if now.Sub(info.ModTime()) <= time.Hour {
			continue
		}
		if err := os.Remove(filepath.Join(dir, entry.Name())); err != nil && !os.IsNotExist(err) {
			return err
		}
		delete(s.pendingDocuments, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())))
	}
	return nil
}

func (s *Service) sweepBoundDocumentDirectory(ctx context.Context) error {
	root := filepath.Join(s.documentRoot, "store-documents", "stores")
	stores, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, store := range stores {
		if !store.IsDir() || !validDocumentStoreID(store.Name()) {
			continue
		}
		if err := s.sweepOneStoreDirectory(ctx, root, store.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) sweepOneStoreDirectory(ctx context.Context, root, storeID string) error {
	dir := filepath.Join(root, storeID)
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || !validDocumentFilename(file.Name()) {
			continue
		}
		url := "/uploads/stores/" + storeID + "/" + file.Name()
		var count int64
		err := s.db.WithContext(ctx).Model(&gen.Store{}).
			Where("id = ? AND (is_delete IS NULL OR is_delete = 1) AND (business_license_image_url = ? OR other_document_image_url = ?)", storeID, url, url).Count(&count).Error
		if err != nil {
			return err
		}
		if count == 0 {
			if err := os.Remove(filepath.Join(dir, file.Name())); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	_ = os.Remove(dir)
	return nil
}
