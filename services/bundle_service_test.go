package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteBrewfileWritesResolvedPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Brewfile")
	service := NewBundleService(nil)
	if err := service.WriteBrewfile(path, "brew \"wget\"\n"); err != nil {
		t.Fatalf("WriteBrewfile returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(data) != "brew \"wget\"\n" {
		t.Fatalf("unexpected Brewfile content: %q", string(data))
	}
}
