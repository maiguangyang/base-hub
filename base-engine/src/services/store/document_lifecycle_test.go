package store

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"base-engine/auth"
	"base-engine/gen"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func documentTestService(t *testing.T) (*Service, *auth.WorkspacePrincipal) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, IgnoreRelationshipsWhenMigrating: true})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&gen.Store{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&gen.Store{ID: "store-1", Code: "S1", Name: "Store", OrganizationID: "org-1", Lifecycle: gen.StoreLifecycleDraft}).Error; err != nil {
		t.Fatal(err)
	}
	org := "org-1"
	p := &auth.WorkspacePrincipal{AccountID: "account-1", SessionID: "session-1", WorkspaceType: auth.WorkspaceTypeFranchise, OrganizationID: &org, AllStores: true, Permissions: map[string]struct{}{"store:create": {}, "store:update": {}, "store:read": {}}}
	s := NewService(db, nil)
	s.SetDocumentRoot(t.TempDir())
	return s, p
}

func documentPNG(t *testing.T) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := png.Encode(&b, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestStoreDocumentBindReadAndRemove(t *testing.T) {
	s, p := documentTestService(t)
	id, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	item, err := s.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindBusinessLicense, id)
	if err != nil {
		t.Fatal(err)
	}
	if item.BusinessLicenseImageURL == nil {
		t.Fatal("image URL missing")
	}
	file, contentType, err := s.OpenDocument(t.Context(), p, "store-1", (*item.BusinessLicenseImageURL)[len("/uploads/stores/store-1/"):])
	if err != nil || contentType != "image/png" {
		t.Fatalf("read: %v %s", err, contentType)
	}
	_ = file.Close()
	if _, err := s.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindOther, id); err == nil {
		t.Fatal("attachment reused")
	}
	removed, err := s.RemoveDocument(t.Context(), p, "store-1", gen.StoreDocumentKindBusinessLicense)
	if err != nil || removed.BusinessLicenseImageURL != nil {
		t.Fatalf("remove: %+v %v", removed, err)
	}
	if _, err := os.Stat(s.storeDocumentPath("store-1", filepath.Base(*item.BusinessLicenseImageURL))); !os.IsNotExist(err) {
		t.Fatalf("removed file remains: %v", err)
	}
}

func TestStoreDocumentReplacementRemovesOldFile(t *testing.T) {
	s, p := documentTestService(t)
	first, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	before, err := s.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindOther, first)
	if err != nil {
		t.Fatal(err)
	}
	oldFile := s.storeDocumentPath("store-1", filepath.Base(*before.OtherDocumentImageURL))
	second, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	after, err := s.SetDocument(t.Context(), p, "store-1", gen.StoreDocumentKindOther, second)
	if err != nil || *after.OtherDocumentImageURL == *before.OtherDocumentImageURL {
		t.Fatalf("replacement: %+v %v", after, err)
	}
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Fatalf("replaced file remains: %v", err)
	}
}

func TestStoreDocumentRejectsForeignSessionAndInvalidImage(t *testing.T) {
	s, p := documentTestService(t)
	if _, err := s.UploadDocument(t.Context(), p, []byte("not an image")); err == nil {
		t.Fatal("invalid image accepted")
	}
	id, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	other := *p
	other.SessionID = "other-session"
	if _, err := s.SetDocument(t.Context(), &other, "store-1", gen.StoreDocumentKindOther, id); err == nil {
		t.Fatal("foreign session bound image")
	}
	if _, _, err := s.OpenDocument(t.Context(), &other, "store-1", "../escape.png"); !os.IsNotExist(err) {
		t.Fatalf("unsafe path: %v", err)
	}
}
