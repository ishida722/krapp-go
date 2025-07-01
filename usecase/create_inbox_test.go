package usecase

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type testInboxConfig struct {
	baseDir       string
	inboxDir      string
	inboxTemplate map[string]any
}

func (c *testInboxConfig) GetBaseDir() string            { return c.baseDir }
func (c *testInboxConfig) GetInboxDir() string           { return c.inboxDir }
func (c *testInboxConfig) GetInboxTemplate() map[string]any { return c.inboxTemplate }

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

func TestCreateInboxNoteWithTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &testInboxConfig{
		baseDir:  tmpDir,
		inboxDir: "inbox",
		inboxTemplate: map[string]any{
			"tags":     []string{"inbox", "test"},
			"status":   "review",
			"priority": "high",
		},
	}
	now := time.Date(2025, 6, 2, 12, 0, 0, 0, time.UTC)
	title := "test-title"
	path, err := CreateInboxNote(cfg, now, title)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// ファイルが作成されたことを確認
	expectedFile := filepath.Join(tmpDir, "inbox", "2025-06-02-test-title.md")
	if path != expectedFile {
		t.Errorf("expected path %s, got %s", expectedFile, path)
	}
	
	// ファイル内容を確認
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, `created: "2025-06-02"`) {
		t.Errorf("created field not found in frontmatter")
	}
	if !strings.Contains(contentStr, "- inbox") && !strings.Contains(contentStr, "- test") {
		t.Errorf("tags field not found in frontmatter")
	}
	if !strings.Contains(contentStr, "status: review") {
		t.Errorf("status field not found in frontmatter")
	}
	if !strings.Contains(contentStr, "priority: high") {
		t.Errorf("priority field not found in frontmatter")
	}
}
