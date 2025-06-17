package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Dummy FrontMatter implementation for testing if not present
type testFrontMatter struct {
	Title   string    `yaml:"title"`
	Created time.Time `yaml:"created"`
}

func (fm *testFrontMatter) SetCreatedNow() {
	fm.Created = time.Now()
}
func (fm *testFrontMatter) ToYAML() (string, error) {
	return "---\ntitle: " + fm.Title + "\ncreated: " + fm.Created.Format(time.RFC3339) + "\n---", nil
}

func TestCreateNewNote(t *testing.T) {
	note, err := CreateNewNote(NewNote{
		Content:   "Hello, world!",
		FilePath:  "",
		WriteFile: false,
		Now:       true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Content != "Hello, world!" {
		t.Errorf("expected content to be 'Hello, world!', got '%s'", note.Content)
	}
}

func TestCreateNewNote_WriteFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "note.md")
	note, err := CreateNewNote(NewNote{
		Content:   "Test file content",
		FilePath:  tmpFile,
		WriteFile: true,
		Now:       false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if !strings.Contains(string(data), "Test file content") {
		t.Errorf("file does not contain expected content")
	}
	if note.FilePath != tmpFile {
		t.Errorf("expected FilePath to be %s, got %s", tmpFile, note.FilePath)
	}
}

func TestCreateNewNote_WriteFile_EmptyPath(t *testing.T) {
	_, err := CreateNewNote(NewNote{
		Content:   "Test",
		FilePath:  "",
		WriteFile: true,
		Now:       false,
	})
	if err == nil {
		t.Error("expected error when FilePath is empty and WriteFile is true")
	}
}

func TestNote_ToString_NoFrontMatter(t *testing.T) {
	note := Note{
		FrontMatter: nil,
		Content:     "Body only",
	}
	s, err := note.ToString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != "Body only" {
		t.Errorf("expected 'Body only', got '%s'", s)
	}
}

func TestNote_ToString_WithFrontMatter(t *testing.T) {
	fm := FrontMatter{}
	if setNow, ok := any(&fm).(interface{ SetCreatedNow() }); ok {
		setNow.SetCreatedNow()
	}
	note := Note{
		FrontMatter: fm,
		Content:     "Body content",
	}
	s, err := note.ToString()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(s, "Body content") {
		t.Errorf("expected output to contain body content")
	}
	if !strings.Contains(s, "---") {
		t.Errorf("expected output to contain frontmatter")
	}
}

func TestParseNote_NoFrontMatter(t *testing.T) {
	raw := "Just content"
	note, err := ParseNote(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Content != "Just content" {
		t.Errorf("expected content to be 'Just content', got '%s'", note.Content)
	}
}

func TestParseNote_WithFrontMatter(t *testing.T) {
	raw := "---\ntitle: test\n---\nBody here"
	note, err := ParseNote(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Content != "Body here" {
		t.Errorf("expected content to be 'Body here', got '%s'", note.Content)
	}
}

func TestParseNote_InvalidFrontMatter(t *testing.T) {
	raw := "---\ntitle: test"
	_, err := ParseNote(raw)
	if err == nil {
		t.Error("expected error for invalid frontmatter format")
	}
}

func TestLoadNoteFromFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "note.md")
	content := "---\ntitle: test\n---\nBody"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	note, err := LoadNoteFromFile(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Content != "Body" {
		t.Errorf("expected content to be 'Body', got '%s'", note.Content)
	}
	if note.FilePath != tmpFile {
		t.Errorf("expected FilePath to be %s, got %s", tmpFile, note.FilePath)
	}
}

func TestLoadNoteFromFile_NotExist(t *testing.T) {
	_, err := LoadNoteFromFile("not_exist_file.md")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestNote_SaveToFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "note.md")
	note := Note{
		FrontMatter: FrontMatter{},
		Content:     "Save this",
		FilePath:    tmpFile,
	}
	err := note.SaveToFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if !strings.Contains(string(data), "Save this") {
		t.Errorf("file does not contain expected content")
	}
}

func TestNote_SaveToFile_EmptyPath(t *testing.T) {
	note := Note{
		FrontMatter: FrontMatter{},
		Content:     "No path",
		FilePath:    "",
	}
	err := note.SaveToFile()
	if err == nil {
		t.Error("expected error when FilePath is empty")
	}
}

func TestNote_MoveFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "move.md")
	if err := os.WriteFile(tmpFile, []byte("move me"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	note := &Note{
		FrontMatter: FrontMatter{},
		Content:     "move me",
		FilePath:    tmpFile,
	}
	newDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(newDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	err := note.MoveFile(newDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedPath := filepath.Join(newDir, "move.md")
	if note.FilePath != expectedPath {
		t.Errorf("expected FilePath to be %s, got %s", expectedPath, note.FilePath)
	}
	if _, err := os.Stat(expectedPath); err != nil {
		t.Errorf("expected file to exist at new location")
	}
}

func TestNote_MoveFile_EmptyPath(t *testing.T) {
	note := &Note{
		FrontMatter: FrontMatter{},
		Content:     "no path",
		FilePath:    "",
	}
	err := note.MoveFile("anywhere")
	if err == nil {
		t.Error("expected error when FilePath is empty")
	}
}
