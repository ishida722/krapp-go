package usecase

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type testDailyConfig struct {
	baseDir       string
	dailyNoteDir  string
	dailyTemplate map[string]any
}

func (c *testDailyConfig) GetBaseDir() string               { return c.baseDir }
func (c *testDailyConfig) GetDailyNoteDir() string          { return c.dailyNoteDir }
func (c *testDailyConfig) GetDailyTemplate() map[string]any { return c.dailyTemplate }

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

func TestCreateDailyNoteWithTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &testDailyConfig{
		baseDir:      tmpDir,
		dailyNoteDir: "daily",
		dailyTemplate: map[string]any{
			"tags":   []string{"daily", "test"},
			"status": "draft",
		},
	}
	now := time.Date(2025, 6, 2, 12, 0, 0, 0, time.UTC)
	path, err := CreateDailyNote(cfg, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ファイルが作成されたことを確認
	expectedFile := filepath.Join(tmpDir, "daily", "2025", "06", "2025-06-02.md")
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
	if !strings.Contains(contentStr, "- daily") && !strings.Contains(contentStr, "- test") {
		t.Errorf("tags field not found in frontmatter")
	}
	if !strings.Contains(contentStr, "status: draft") {
		t.Errorf("status field not found in frontmatter")
	}
}

func TestCreateDailyNoteExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &testDailyConfig{
		baseDir:      tmpDir,
		dailyNoteDir: "daily",
	}
	now := time.Date(2025, 6, 2, 12, 0, 0, 0, time.UTC)
	expectedFile := filepath.Join(tmpDir, "daily", "2025", "06", "2025-06-02.md")
	
	// ディレクトリを作成
	if err := os.MkdirAll(filepath.Dir(expectedFile), 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	
	// 既存のファイルを作成
	existingContent := `---
created: "2025-06-02"
tags: ["existing"]
---

# 既存のノート

これは既存のノートです。`
	
	if err := os.WriteFile(expectedFile, []byte(existingContent), 0644); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}
	
	// CreateDailyNoteを呼び出す
	path, err := CreateDailyNote(cfg, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// パスが正しいことを確認
	if path != expectedFile {
		t.Errorf("expected path %s, got %s", expectedFile, path)
	}
	
	// ファイルが上書きされていないことを確認
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "既存のノート") {
		t.Errorf("existing content was overwritten")
	}
	if !strings.Contains(contentStr, `tags: ["existing"]`) {
		t.Errorf("existing tags were overwritten")
	}
}
