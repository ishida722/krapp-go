package usecase

import (
	"fmt"
	"os/exec"
)

func SyncGit() error {
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return fmt.Errorf("git add に失敗: %w", err)
	}
	if err := exec.Command("git", "commit", "-m", "add").Run(); err != nil {
		return fmt.Errorf("git commit に失敗: %w", err)
	}
	if err := exec.Command("git", "pull").Run(); err != nil {
		return fmt.Errorf("git pull に失敗: %w", err)
	}
	if err := exec.Command("git", "push").Run(); err != nil {
		return fmt.Errorf("git push に失敗: %w", err)
	}
	return nil
}
