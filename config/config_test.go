package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Default(t *testing.T) {
	cfg, err := LoadConfig("not_exist_config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "./notes", cfg.BaseDir)
	assert.Equal(t, "daily", cfg.DailyNoteDir)
}

func TestMarshalYAML(t *testing.T) {
	cfg := &Config{BaseDir: "/tmp", DailyNoteDir: "dailies"}
	b, err := MarshalYAML(cfg)
	assert.NoError(t, err)
	assert.NotEmpty(t, b, "MarshalYAML should return non-empty byte slice")
}

func TestSaveAndLoadConfig(t *testing.T) {
	tmpFile := "test_config.yaml"
	defer os.Remove(tmpFile)
	cfg := &Config{BaseDir: "/test", DailyNoteDir: "testdaily"}
	// ここでセーブ
	err := SaveConfig(tmpFile, cfg)
	assert.NoError(t, err, "SaveConfig should not return an error")
	// セーブしたものをロード
	loaded, err := LoadConfig(tmpFile)
	assert.NoError(t, err, "LoadConfig should not return an error")
	assert.Equal(t, cfg.BaseDir, loaded.BaseDir, "BaseDir should match")
	assert.Equal(t, cfg.DailyNoteDir, loaded.DailyNoteDir, "DailyNoteDir should match")
}

func TestMergeConfig(t *testing.T) {
	cfg1 := &Config{BaseDir: "/base1", DailyNoteDir: "daily1"}
	cfg2 := &Config{BaseDir: "/base2", Inbox: "inbox2"}
	merged := MergeConfig(*cfg1, *cfg2)
	assert.Equal(t, "/base1", merged.BaseDir, "BaseDir should be from cfg1")
	assert.Equal(t, "daily1", merged.DailyNoteDir, "DailyNoteDir should be from cfg1")

}
