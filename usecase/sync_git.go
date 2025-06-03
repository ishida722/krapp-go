package usecase

import (
	"fmt"
	"os"
	"os/exec"
)

func SyncGit(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("指定されたディレクトリが存在しません: %s", dir)
	}
	fmt.Println("sync:", dir)
	cmds := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "add"},
		{"git", "pull"},
		{"git", "push"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir // ここで実行ディレクトリを指定
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s に失敗: %w", args[0], err)
		}
	}
	return nil
}
