package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveCreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "shots", "page.png")

	f := New()
	if err := f.Save(path, []byte("png-data")); err != nil {
		t.Fatalf("save: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved file: %v", err)
	}
	if string(data) != "png-data" {
		t.Fatalf("content: got %q, want %q", data, "png-data")
	}
}

func TestSaveToExistingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "page.png")

	f := New()
	if err := f.Save(path, []byte("data")); err != nil {
		t.Fatalf("save: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not found: %v", err)
	}
}