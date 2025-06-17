package usecase

import (
	"testing"

	"github.com/ishida722/krapp-go/models"
)

func TestAddCreatedFromNote_FileName(t *testing.T) {
	note := &models.Note{
		Content:     "some content",
		FilePath:    "/notes/2025-06-10-title.md",
		FrontMatter: models.FrontMatter{},
	}
	if err := AddCreatedFromNote(note); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.FrontMatter["Created"] != "2025-06-10" {
		t.Errorf("expected Created 2025-06-10, got %v", note.FrontMatter["Created"])
	}
}

func TestAddCreatedFromNote_Content(t *testing.T) {
	note := &models.Note{
		Content:     "log from 2025-06-11 about something",
		FilePath:    "note.md",
		FrontMatter: models.FrontMatter{},
	}
	if err := AddCreatedFromNote(note); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.FrontMatter["Created"] != "2025-06-11" {
		t.Errorf("expected Created 2025-06-11, got %v", note.FrontMatter["Created"])
	}
}

func TestAddCreatedFromNote_NoDate(t *testing.T) {
	note := &models.Note{Content: "no date here"}
	if err := AddCreatedFromNote(note); err == nil {
		t.Error("expected error, got nil")
	}
}
