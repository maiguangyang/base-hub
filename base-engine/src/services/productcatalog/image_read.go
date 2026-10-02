package productcatalog

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"base-engine/gen"
)

func (s *Service) OpenMainImage(ctx context.Context, productID, filename string) (*os.File, string, error) {
	url := "/uploads/products/" + productID + "/" + filename
	path, ok := managedImagePath(s.imageRoot, productID, url)
	if !ok {
		return nil, "", os.ErrNotExist
	}
	var count int64
	if err := s.db.WithContext(ctx).Model(&gen.Product{}).Where("id = ? AND image_url = ?", productID, url).Count(&count).Error; err != nil {
		return nil, "", err
	}
	if count != 1 {
		return nil, "", os.ErrNotExist
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	contentType := "image/jpeg"
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".png":
		contentType = "image/png"
	case ".webp":
		contentType = "image/webp"
	}
	return file, contentType, nil
}
