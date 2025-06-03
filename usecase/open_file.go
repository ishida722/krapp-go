package usecase

import (
	"fmt"
	"os"
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
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ファイルを開く際にエラー: %w", err)
	}
	return nil
}
