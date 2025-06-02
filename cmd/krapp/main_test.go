package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) {
	notesDir := filepath.Join("..", "..", "notes")
	_ = os.RemoveAll(notesDir)
}

func TestPrintConfig(t *testing.T) {
	setup(t)
	cmd := exec.Command("go", "run", "main.go", "print-config")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("print-config failed: %v\nOutput: %s", err, string(out))
	}
	if len(out) == 0 {
		t.Error("print-config output is empty")
	}
}

func TestCreateDaily(t *testing.T) {
	setup(t)
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

func TestCreateInbox(t *testing.T) {
	setup(t)
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
