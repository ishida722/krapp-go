package usecase

import (
	"fmt"
	"os/exec"
)

type OpenFileConfig interface {
	GetEditorCommand() string
}

func OpenFile(editorCommand, filePath string) error {
	if editorCommand == "" {
		// エディタコマンドが指定されていない場合はエラー
		return fmt.Errorf("エディタが指定されていません")
	}

	cmd := exec.Command(editorCommand, filePath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ファイルを開く際にエラー: %w", err)
	}
	return nil
}
