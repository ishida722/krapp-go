package usecase

import (
	"fmt"
	"os"

	"github.com/ishida722/krapp-go/usecase/frontmatter"
)

func addCreatedForNewFile(filePath string) error {
	content, err := frontmatter.AddCreated("")
	if err != nil {
		return fmt.Errorf("frontmatterの追加に失敗: %w", err)
	}
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("ファイル書き込みに失敗: %w", err)
	}
	return nil
}
