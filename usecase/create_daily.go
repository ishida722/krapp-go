package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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
	filePath := filepath.Join(dir, date+".md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		f, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("ファイル作成に失敗: %w", err)
		}
		f.Close()
	}
	return filePath, nil
}
