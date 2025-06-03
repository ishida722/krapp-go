package usecase

import (
	"fmt"
	"os"
	"os/exec"
)

func SyncGit() error {
	cmds := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "add"},
		{"git", "pull"},
		{"git", "push"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s に失敗: %w", args[0], err)
		}
	}
	return nil
}
