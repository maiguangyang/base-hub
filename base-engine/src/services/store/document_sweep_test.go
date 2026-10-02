package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreDocumentSweepRemovesExpiredPendingFile(t *testing.T) {
	s, p := documentTestService(t)
	id, err := s.UploadDocument(t.Context(), p, documentPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	path := s.pendingDocumentPath(id, "png")
	old := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
	if err := s.SweepDocumentFiles(t.Context()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expired pending file remains: %v", err)
	}
	if _, ok := s.pendingDocuments[id]; ok {
		t.Fatal("expired pending token remains")
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("pending directory missing: %v", err)
	}
}
