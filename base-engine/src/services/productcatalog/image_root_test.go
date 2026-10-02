package productcatalog

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestDefaultImageRootStaysUnderEngineWhenWorkingDirectoryChanges(t *testing.T) {
	t.Setenv("PRODUCT_IMAGE_UPLOAD_DIR", "")
	t.Chdir(t.TempDir())
	_, source, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("source path unavailable")
	}
	want := filepath.Clean(filepath.Join(filepath.Dir(source), "..", "..", "..", "uploads"))
	if got := NewService(nil, nil).ImageRoot(); got != want {
		t.Fatalf("image root = %q, want %q", got, want)
	}
}
