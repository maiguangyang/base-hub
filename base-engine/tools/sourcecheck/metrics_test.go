package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMeasureFileExcludesCommentsAndCountsControlFlow 验证有效行和控制流度量。
func TestMeasureFileExcludesCommentsAndCountsControlFlow(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.go")
	source := "package sample\n// comment\nfunc small(value bool) bool {\n\tif value {\n\t\tfor value {\n\t\t\treturn false\n\t\t}\n\t}\n\treturn true\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	items, err := measureFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := items[0].Metrics.Lines; got != 9 {
		t.Fatalf("file active lines = %d, want 9", got)
	}
	if got := items[1].Metrics; got.Lines != 8 || got.Complexity != 3 || got.Depth != 2 {
		t.Fatalf("function metrics = %+v", got)
	}
}

// TestSourceFilesNeverReadsGeneratedDirectory 验证生成目录不会进入扫描列表。
func TestSourceFilesNeverReadsGeneratedDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "gen"), 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"main.go", "gen/generated.go", "resolver_gen.go"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("package sample"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	files, err := sourceFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || !strings.HasSuffix(files[0], "main.go") {
		t.Fatalf("source files = %v", files)
	}
}

// TestLegacyBaselineDoesNotGrow 验证既有超限函数不能继续增长。
func TestLegacyBaselineDoesNotGrow(t *testing.T) {
	item := measurement{Name: "old", Metrics: metrics{Lines: 51, Complexity: 11, Depth: 3}}
	baseline := map[string]metrics{"old.go#old": {Lines: 50, Complexity: 11, Depth: 3}}
	if !checkMeasurement("old.go", item, baseline, false) {
		t.Fatal("a line beyond the recorded baseline must fail")
	}
}

// TestMeasureFileIncludesClosures 验证匿名函数也参与结构限制。
func TestMeasureFileIncludesClosures(t *testing.T) {
	path := filepath.Join(t.TempDir(), "closure.go")
	source := "package sample\nvar choose = func(value bool) bool {\n if value { return true }; return false\n}\n"
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	items, err := measureFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[1].Name != "closure@2" || items[1].Metrics.Complexity != 2 {
		t.Fatalf("closure metrics = %+v", items)
	}
}
