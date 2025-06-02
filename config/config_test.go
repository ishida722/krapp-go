package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Default(t *testing.T) {
	cfg, err := LoadConfig("not_exist_config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseDir != "./notes" {
		t.Errorf("expected BaseDir './notes', got '%s'", cfg.BaseDir)
	}
	if cfg.DailyNoteDir != "daily" {
		t.Errorf("expected DailyNoteDir 'daily', got '%s'", cfg.DailyNoteDir)
	}
}

func TestMarshalYAML(t *testing.T) {
	cfg := &Config{BaseDir: "/tmp", DailyNoteDir: "dailies"}
	b, err := MarshalYAML(cfg)
	if err != nil {
		t.Fatalf("MarshalYAML error: %v", err)
	}
	if string(b) == "" {
		t.Error("MarshalYAML returned empty string")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpFile := "test_config.yaml"
	defer os.Remove(tmpFile)
	cfg := &Config{BaseDir: "/test", DailyNoteDir: "testdaily"}
	err := SaveConfig(tmpFile, cfg)
	if err != nil {
		t.Fatalf("SaveConfig error: %v", err)
	}
	loaded, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if loaded.BaseDir != cfg.BaseDir || loaded.DailyNoteDir != cfg.DailyNoteDir {
		t.Errorf("Loaded config does not match saved config")
	}
}
