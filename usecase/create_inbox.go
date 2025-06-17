package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ishida722/krapp-go/models"
)

type InboxConfig interface {
	GetBaseDir() string
	GetInboxDir() string
}

// CreateInboxNote creates a new inbox note with the given title and returns its path.
func CreateInboxNote(cfg InboxConfig, now time.Time, title string) (string, error) {
	date := now.Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.md", date, title)
	dir := filepath.Join(cfg.GetBaseDir(), cfg.GetInboxDir())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("inboxディレクトリ作成に失敗: %w", err)
	}
	note, err := models.CreateNewNote(models.NewNote{
		Content:   "",
		FilePath:  filepath.Join(dir, filename),
		WriteFile: true,
		Now:       true,
	})
	if err != nil {
		return "", fmt.Errorf("日記の保存に失敗: %w", err)
	}
	return note.FilePath, nil
}
