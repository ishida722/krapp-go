package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ishida722/krapp-go/usecase/frontmatter"
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
	filePath := filepath.Join(dir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("inboxノート作成に失敗: %w", err)
		}
		f.Close()
	}
	content, err := frontmatter.AddCreated("")
	if err != nil {
		return "", fmt.Errorf("frontmatterの追加に失敗: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("ファイル書き込みに失敗: %w", err)
	}
	return filePath, nil
}
