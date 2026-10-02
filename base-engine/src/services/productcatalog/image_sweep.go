package productcatalog

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"base-engine/gen"
)

func (s *Service) SweepImageFiles(ctx context.Context) error {
	s.imageMu.Lock()
	defer s.imageMu.Unlock()
	now := time.Now()
	s.sweepPendingImages(now)
	err := s.sweepDiskImages(ctx, now)
	if err == nil {
		s.lastImageSweep = now
	} else {
		s.lastImageSweep = time.Time{}
	}
	return err
}

func (s *Service) sweepDiskImages(ctx context.Context, now time.Time) error {
	if err := s.sweepPendingDirectory(now); err != nil {
		return err
	}
	return s.sweepProductDirectory(ctx)
}

func (s *Service) sweepPendingDirectory(now time.Time) error {
	directory := filepath.Join(s.imageRoot, "pending")
	entries, err := os.ReadDir(directory)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.IsDir() || now.Sub(info.ModTime()) <= time.Hour {
			continue
		}
		if err := os.Remove(filepath.Join(directory, entry.Name())); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (s *Service) sweepProductDirectory(ctx context.Context) error {
	directory := filepath.Join(s.imageRoot, "products")
	products, err := os.ReadDir(directory)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, product := range products {
		if !product.IsDir() {
			continue
		}
		if err := s.sweepProductImages(ctx, product.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) sweepProductImages(ctx context.Context, productID string) error {
	directory := filepath.Join(s.imageRoot, "products", productID)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		url := "/uploads/products/" + productID + "/" + entry.Name()
		if _, valid := managedImagePath(s.imageRoot, productID, url); !valid {
			continue
		}
		var count int64
		if err := s.db.WithContext(ctx).Model(&gen.Product{}).Where("id = ? AND image_url = ?", productID, url).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := os.Remove(filepath.Join(directory, entry.Name())); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	_ = os.Remove(directory)
	return nil
}
