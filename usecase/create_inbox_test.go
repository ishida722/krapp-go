package usecase

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testInboxConfig struct {
	baseDir  string
	inboxDir string
}

func (c *testInboxConfig) GetBaseDir() string  { return c.baseDir }
func (c *testInboxConfig) GetInboxDir() string { return c.inboxDir }

func TestCreateInboxNote(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &testInboxConfig{
		baseDir:  tmpDir,
		inboxDir: "inbox",
	}
	now := time.Date(2025, 6, 2, 12, 0, 0, 0, time.UTC)
	title := "test-title"
	path, err := CreateInboxNote(cfg, now, title)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedFile := filepath.Join(tmpDir, "inbox", "2025-06-02-test-title.md")
	if path != expectedFile {
		t.Errorf("expected path %s, got %s", expectedFile, path)
	}
	if _, err := os.Stat(expectedFile); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
