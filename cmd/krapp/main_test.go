package main

import (
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestPrintConfig(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "config")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("config failed: %v\nOutput: %s", err, string(out))
	}
	if len(out) == 0 {
		t.Error("config output is empty")
	}
}

func TestCreateDaily(t *testing.T) {
	// テスト用の設定ファイルを一時的にコピー
	if err := copyTestConfig(); err != nil {
		t.Fatalf("Failed to setup test config: %v", err)
	}
	defer cleanupTestConfig()

	cmd := exec.Command("go", "run", "main.go", "create-daily")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("create-daily failed: %v\nOutput: %s", err, string(out))
	}
	if len(out) == 0 {
		t.Error("create-daily output is empty")
	}
}

func TestCreateDailyShort(t *testing.T) {
	if err := copyTestConfig(); err != nil {
		t.Fatalf("Failed to setup test config: %v", err)
	}
	defer cleanupTestConfig()

	cmd := exec.Command("go", "run", "main.go", "cd")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("create-daily failed: %v\nOutput: %s", err, string(out))
	}
	if len(out) == 0 {
		t.Error("create-daily output is empty")
	}
}

func TestCreateInbox(t *testing.T) {
	if err := copyTestConfig(); err != nil {
		t.Fatalf("Failed to setup test config: %v", err)
	}
	defer cleanupTestConfig()

	cmd := exec.Command("go", "run", "main.go", "create-inbox", "test")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("create-inbox failed: %v\nOutput: %s", err, string(out))
	}
	if len(out) == 0 {
		t.Error("create-inbox output is empty")
	}
}

func TestCreateInboxShort(t *testing.T) {
	if err := copyTestConfig(); err != nil {
		t.Fatalf("Failed to setup test config: %v", err)
	}
	defer cleanupTestConfig()

	cmd := exec.Command("go", "run", "main.go", "ci", "test")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("create-inbox failed: %v\nOutput: %s", err, string(out))
	}
	if len(out) == 0 {
		t.Error("create-inbox output is empty")
	}
}

func copyTestConfig() error {
	src, err := os.Open("test_config.yaml")
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(".krapp_config.yaml.backup")
	if err != nil {
		return err
	}
	defer dst.Close()

	// 現在の設定をバックアップ
	if _, err := os.Stat(".krapp_config.yaml"); err == nil {
		existing, err := os.Open(".krapp_config.yaml")
		if err != nil {
			return err
		}
		defer existing.Close()
		if _, err := io.Copy(dst, existing); err != nil {
			return err
		}
	}

	// テスト用設定をコピー
	dst2, err := os.Create(".krapp_config.yaml")
	if err != nil {
		return err
	}
	defer dst2.Close()

	if _, err := io.Copy(dst2, src); err != nil {
		return err
	}

	return nil
}

func cleanupTestConfig() error {
	// バックアップを復元
	if _, err := os.Stat(".krapp_config.yaml.backup"); err == nil {
		if err := os.Rename(".krapp_config.yaml.backup", ".krapp_config.yaml"); err != nil {
			return err
		}
	} else {
		// バックアップがない場合は削除
		os.Remove(".krapp_config.yaml")
	}
	return nil
}
