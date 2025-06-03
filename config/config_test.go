package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup_config_file(t *testing.T) string {
	tempDir := t.TempDir()
	cfgGlobal := GetDefaultConfig()
	cfgGlobal.BaseDir = tempDir
	cfgLocal := cfgGlobal
	// エディタだけローカル設定を変更
	cfgLocal.Editor = "nvim"
	SetConfigPaths(ConfigPaths{
		Global: filepath.Join(tempDir, ".krapp_config.yaml"),
		Local:  filepath.Join(tempDir, ".krapp_config_local.yaml"),
	})
	cfgPaths, _ := GetConfigPaths()
	saveConfig(cfgPaths.Global, cfgGlobal)
	saveConfig(cfgPaths.Local, cfgLocal)
	return tempDir
}

func TestLoadConfig_Default(t *testing.T) {
	cfg, err := loadConfig("not_exist_config.yaml")
	assert.Error(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.BaseDir)
	assert.Equal(t, "", cfg.DailyNoteDir)
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
	cfg := Config{BaseDir: "/test", DailyNoteDir: "testdaily"}
	// ここでセーブ
	err := saveConfig(tmpFile, cfg)
	assert.NoError(t, err, "SaveConfig should not return an error")
	// セーブしたものをロード
	loaded, err := loadConfig(tmpFile)
	assert.NoError(t, err, "LoadConfig should not return an error")
	assert.Equal(t, cfg.BaseDir, loaded.BaseDir, "BaseDir should match")
	assert.Equal(t, cfg.DailyNoteDir, loaded.DailyNoteDir, "DailyNoteDir should match")
}

func TestNullValueConfig(t *testing.T) {
	// 値を指定しないと""が設定される"
	cfg1 := &Config{BaseDir: "/base1", DailyNoteDir: "daily1"}
	assert.Equal(t, "/base1", cfg1.BaseDir)
	assert.Equal(t, "daily1", cfg1.DailyNoteDir)
	assert.Equal(t, "", cfg1.Editor)
	assert.Equal(t, "", cfg1.Inbox)

	cfg2 := &Config{BaseDir: "/base2", Inbox: "inbox2"}
	assert.Equal(t, "/base2", cfg2.BaseDir)
	assert.Equal(t, "", cfg2.DailyNoteDir)
	assert.Equal(t, "", cfg2.Editor)
	assert.Equal(t, "inbox2", cfg2.Inbox)
}

func TestMergeConfig(t *testing.T) {
	// グローバル設定の一部をローカル設定で上書きする
	global_cfg := Config{
		BaseDir:      "/base",
		DailyNoteDir: "daily",
		Inbox:        "inbox",
		Editor:       "vim",
	}
	// ローカル設定でEditorを上書き
	local_cfg := Config{Editor: "code"}
	merged := MergeConfig(global_cfg, local_cfg)

	assert.Equal(t, "/base", merged.BaseDir)
	assert.Equal(t, "daily", merged.DailyNoteDir)
	assert.Equal(t, "inbox", merged.Inbox)
	assert.Equal(t, "code", merged.Editor, "Editor should be overwritten by local config")

}

func TestLoadMergedConfig(t *testing.T) {
	tempDir := setup_config_file(t)

	cfg, err := LoadConfig()
	assert.NoError(t, err, "LoadConfig should not return an error")
	assert.Equal(t, cfg.BaseDir, tempDir, "BaseDir should match the temp directory")
	assert.Equal(t, cfg.DailyNoteDir, "daily", "DailyNoteDir should match default value")
	assert.Equal(t, cfg.Inbox, "inbox", "Inbox should match default value")
	assert.Equal(t, cfg.Editor, "nvim", "Editor should be overridden by local config")
}
