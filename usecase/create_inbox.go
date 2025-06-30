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
	GetInboxTemplate() map[string]any
}

// CreateInboxNote creates a new inbox note with the given title and returns its path.
func CreateInboxNote(cfg InboxConfig, now time.Time, title string) (string, error) {
	date := now.Format("2006-01-02")
	filename := fmt.Sprintf("%s-%s.md", date, title)
	dir := filepath.Join(cfg.GetBaseDir(), cfg.GetInboxDir())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("inboxディレクトリ作成に失敗: %w", err)
	}

	// テンプレートから初期frontmatterを作成
	fm := models.FrontMatter{}
	
	// 作成日時を設定
	fm.SetCreated(now)
	
	// テンプレートの属性を追加
	if cfg.GetInboxTemplate() != nil {
		for key, value := range cfg.GetInboxTemplate() {
			fm[key] = value
		}
	}

	note, err := models.CreateNewNoteWithFrontMatter(models.NewNoteWithFrontMatter{
		Content:     "",
		FilePath:    filepath.Join(dir, filename),
		WriteFile:   true,
		FrontMatter: fm,
	})
	if err != nil {
		return "", fmt.Errorf("日記の保存に失敗: %w", err)
	}
	return note.FilePath, nil
}
