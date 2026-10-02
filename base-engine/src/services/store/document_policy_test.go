package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"base-engine/auth"
	"base-engine/gen"
)

func TestStoreDocumentRejectsOversizeAndExpiredAttachment(t *testing.T) {
	s, p := documentTestService(t)
	if _, err := s.UploadDocument(t.Context(), p, make([]byte, maxStoreDocumentBytes+1)); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("oversize upload: %v", err)
	}
	id, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	pending := s.pendingDocuments[id]
	pending.createdAt = time.Now().Add(-2 * time.Hour)
	s.pendingDocuments[id] = pending
	if _, err := s.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindOther, id); auth.ErrorCode(err) != auth.CodeValidationFailed {
		t.Fatalf("expired attachment: %v", err)
	}
}

func TestStoreDocumentRejectsForeignStoreAndInactiveWrite(t *testing.T) {
	s, p := documentTestService(t)
	id, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	foreign := *p
	org := "other-org"
	foreign.OrganizationID = &org
	if _, err := s.SetDocument(t.Context(), &foreign, "store-1", gen.StoreDocumentKindOther, id); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("foreign bind: %v", err)
	}
	if err := s.db.Model(&gen.Store{}).Where("id = ?", "store-1").Update("lifecycle", gen.StoreLifecycleActive).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindOther, id); auth.ErrorCode(err) != auth.CodeConflict {
		t.Fatalf("active franchise write: %v", err)
	}
}

func TestStoreDocumentPreviewRequiresCurrentStoreAccess(t *testing.T) {
	s, p := documentTestService(t)
	id, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	item, err := s.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindOther, id)
	if err != nil {
		t.Fatal(err)
	}
	filename := filepath.Base(*item.OtherDocumentImageURL)
	if _, _, err := s.OpenDocument(t.Context(), nil, "store-1", filename); !os.IsNotExist(err) {
		t.Fatalf("anonymous preview: %v", err)
	}
	withoutRead := *p
	withoutRead.Permissions = map[string]struct{}{"store:update": {}}
	if _, _, err := s.OpenDocument(t.Context(), &withoutRead, "store-1", filename); !os.IsNotExist(err) {
		t.Fatalf("missing read permission: %v", err)
	}
	wrongOrganization := *p
	foreignOrg := "other-org"
	wrongOrganization.OrganizationID = &foreignOrg
	if _, _, err := s.OpenDocument(t.Context(), &wrongOrganization, "store-1", filename); !os.IsNotExist(err) {
		t.Fatalf("foreign preview: %v", err)
	}
	deleted := 2
	if err := s.db.Model(&gen.Store{}).Where("id = ?", "store-1").Update("is_delete", deleted).Error; err != nil {
		t.Fatal(err)
	}
	if _, _, err := s.OpenDocument(t.Context(), p, "store-1", filename); !os.IsNotExist(err) {
		t.Fatalf("soft-deleted preview: %v", err)
	}
}

func TestStoreDocumentCreateOnlyPermissionBindsOnlyNewOwnStore(t *testing.T) {
	s, p := documentTestService(t)
	createOnly := *p
	createOnly.Permissions = map[string]struct{}{"store:create": {}}
	id, err := s.UploadDocument(t.Context(), &createOnly, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetDocument(t.Context(), &createOnly, "store-1", gen.StoreDocumentKindOther, id); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("old store accepted create-only bind: %v", err)
	}
	createdAt := s.pendingDocuments[id].createdAt.UnixMilli()
	if err := s.db.Model(&gen.Store{}).Where("id = ?", "store-1").Updates(map[string]any{"created_by": p.AccountID, "created_at": createdAt}).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetDocument(t.Context(), &createOnly, "store-1", gen.StoreDocumentKindOther, id); err != nil {
		t.Fatalf("new own store rejected: %v", err)
	}
	if _, err := s.RemoveDocument(t.Context(), &createOnly, "store-1", gen.StoreDocumentKindOther); auth.ErrorCode(err) != auth.CodePermissionDenied {
		t.Fatalf("create-only removal accepted: %v", err)
	}
}

func TestStoreDocumentHQCreateOnlyBindsDirectStore(t *testing.T) {
	s, p := documentTestService(t)
	hq := *p
	hq.WorkspaceType = auth.WorkspaceTypeHeadquarters
	hq.Permissions = map[string]struct{}{"hqStore:create": {}}
	id, err := s.UploadDocument(t.Context(), &hq, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	createdAt := s.pendingDocuments[id].createdAt.UnixMilli()
	if err := s.db.Model(&gen.Store{}).Where("id = ?", "store-1").Updates(map[string]any{
		"created_by": p.AccountID, "created_at": createdAt, "lifecycle": gen.StoreLifecycleActive,
	}).Error; err != nil {
		t.Fatal(err)
	}
	if _, err := s.SetDocument(t.Context(), &hq, "store-1", gen.StoreDocumentKindBusinessLicense, id); err != nil {
		t.Fatalf("HQ create-only bind: %v", err)
	}
}
