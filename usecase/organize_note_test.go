package usecase

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ishida722/krapp-go/models"
)

func createTempNote(t *testing.T, dir, name string, fm models.FrontMatter) models.Note {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return models.Note{FrontMatter: fm, FilePath: path}
}

func TestOrganizeNotesByCreated(t *testing.T) {
	baseDir := t.TempDir()
	// prepare notes
	note1 := createTempNote(t, baseDir, "note1.md", models.FrontMatter{"created": time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)})
	note2 := createTempNote(t, baseDir, "note2.md", models.FrontMatter{"created": time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)})

	// create destination directories
	os.MkdirAll(filepath.Join(baseDir, "2025", "06"), 0755)
	os.MkdirAll(filepath.Join(baseDir, "2025", "07"), 0755)

	OrganizeNotesByCreated([]models.Note{note1, note2}, baseDir)

	expected1 := filepath.Join(baseDir, "2025", "06", "note1.md")
	if _, err := os.Stat(expected1); err != nil {
		t.Fatalf("note1 not moved: %v", err)
	}
	expected2 := filepath.Join(baseDir, "2025", "07", "note2.md")
	if _, err := os.Stat(expected2); err != nil {
		t.Fatalf("note2 not moved: %v", err)
	}
}

func TestOrganizeNotesByLabel(t *testing.T) {
	baseDir := t.TempDir()
	labelMap := LabelDirectoryMap{"diary": "diary", "note": "notes"}

	note1 := createTempNote(t, baseDir, "n1.md", models.FrontMatter{"label": "diary"})
	note2 := createTempNote(t, baseDir, "n2.md", models.FrontMatter{"label": "note"})
	note3 := createTempNote(t, baseDir, "n3.md", models.FrontMatter{"label": "other"})

	os.MkdirAll(filepath.Join(baseDir, "diary"), 0755)
	os.MkdirAll(filepath.Join(baseDir, "notes"), 0755)

	OrganizeNotesByLabel([]models.Note{note1, note2, note3}, baseDir, labelMap)

	if _, err := os.Stat(filepath.Join(baseDir, "diary", "n1.md")); err != nil {
		t.Fatalf("note1 not moved: %v", err)
	}
	if _, err := os.Stat(filepath.Join(baseDir, "notes", "n2.md")); err != nil {
		t.Fatalf("note2 not moved: %v", err)
	}
	// note3 label not in map -> should remain in original location
	if _, err := os.Stat(filepath.Join(baseDir, "n3.md")); err != nil {
		t.Fatalf("note3 should not be moved: %v", err)
	}
}
