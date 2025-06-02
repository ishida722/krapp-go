package usecase

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testDailyConfig struct {
	baseDir      string
	dailyNoteDir string
}

func (c *testDailyConfig) GetBaseDir() string      { return c.baseDir }
func (c *testDailyConfig) GetDailyNoteDir() string { return c.dailyNoteDir }

func TestCreateDailyNote(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &testDailyConfig{
		baseDir:      tmpDir,
		dailyNoteDir: "daily",
	}
	now := time.Date(2025, 6, 2, 12, 0, 0, 0, time.UTC)
	path, err := CreateDailyNote(cfg, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedFile := filepath.Join(tmpDir, "daily", "2025", "06", "2025-06-02.md")
	if path != expectedFile {
		t.Errorf("expected path %s, got %s", expectedFile, path)
	}
	if _, err := os.Stat(expectedFile); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
