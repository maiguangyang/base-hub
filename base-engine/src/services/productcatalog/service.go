package productcatalog

import (
	"base-engine/src/services/audit"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Service owns the headquarters product catalog.
type Service struct {
	db             *gorm.DB
	audit          *audit.Service
	imageRoot      string
	imageMu        sync.Mutex
	lastImageSweep time.Time
	pendingImages  map[string]pendingImage
}

func NewService(db *gorm.DB, auditService *audit.Service) *Service {
	return &Service{db: db, audit: auditService, imageRoot: defaultImageRoot(), pendingImages: make(map[string]pendingImage)}
}

func defaultImageRoot() string {
	if configured := os.Getenv("PRODUCT_IMAGE_UPLOAD_DIR"); configured != "" {
		if root, err := filepath.Abs(configured); err == nil {
			return root
		}
	}
	if cwd, err := os.Getwd(); err == nil && engineModuleExists(cwd) {
		return filepath.Join(cwd, "uploads")
	}
	if _, source, _, ok := runtime.Caller(0); ok && filepath.IsAbs(source) {
		root := filepath.Clean(filepath.Join(filepath.Dir(source), "..", "..", ".."))
		if engineModuleExists(root) {
			return filepath.Join(root, "uploads")
		}
	}
	if executable, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(executable), "uploads")
	}
	return "uploads"
}

func engineModuleExists(root string) bool {
	_, err := os.Stat(filepath.Join(root, "model", "extend.graphql"))
	return err == nil
}
