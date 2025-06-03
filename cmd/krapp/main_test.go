package main

import (
	"os/exec"
	"testing"
)

func TestPrintConfig(t *testing.T) {
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
