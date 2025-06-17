package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ishida722/krapp-go/models"
)

type Config interface {
	GetBaseDir() string
	GetDailyNoteDir() string
}

// CreateDailyNote creates today's daily note and returns its path.
func CreateDailyNote(cfg Config, now time.Time) (string, error) {
	year := now.Format("2006")
	month := now.Format("01")
	date := now.Format("2006-01-02")
	dir := filepath.Join(cfg.GetBaseDir(), cfg.GetDailyNoteDir(), year, month)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("ディレクトリ作成に失敗: %w", err)
	}
	note, err := models.CreateNewNote(models.NewNote{
		Content:   "",
		FilePath:  filepath.Join(dir, date+".md"),
		WriteFile: true,
		Now:       true,
	})
	if err != nil {
		return "", fmt.Errorf("日記の保存に失敗: %w", err)
	}
	return note.FilePath, nil
}
