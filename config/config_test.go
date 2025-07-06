package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupConfigFile(t *testing.T) string {
	tempDir := t.TempDir()
	cfgDefault := GetDefaultConfig()
	cfgDefault.BaseDir = tempDir
	cfgGlobal := cfgDefault
	cfgLocal := cfgGlobal
	// エディタだけローカル設定を変更
	cfgLocal.Editor = "nvim"
	// editorerのオプションも設定
	cfgGlobal.EditorOption = "-c"
	SetConfigPaths(ConfigPaths{
		Global: filepath.Join(tempDir, "config", "krapp", "config.yaml"),
		Local:  filepath.Join(tempDir, ".krapp_config.yaml"),
	})
	cfgPaths, _ := GetConfigPaths()
	// Create directory for global config
	os.MkdirAll(filepath.Dir(cfgPaths.Global), 0755)
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
	tempDir := setupConfigFile(t)

	cfg, err := LoadConfig()
	assert.NoError(t, err, "LoadConfig should not return an error")
	assert.Equal(t, cfg.BaseDir, tempDir, "BaseDir should match the temp directory")
	assert.Equal(t, cfg.DailyNoteDir, "daily", "DailyNoteDir should match default value")
	assert.Equal(t, cfg.Inbox, "inbox", "Inbox should match default value")
	assert.Equal(t, cfg.Editor, "nvim", "Editor should be overridden by local config")

	t.Cleanup(func() {
		// テスト後に設定ファイルを削除
		os.Remove(filepath.Join(tempDir, "config", "krapp", "config.yaml"))
		os.Remove(filepath.Join(tempDir, ".krapp_config.yaml"))
		// 設定パスをリセット
		ResetConfigPaths()
	})
}

// TestXDGConfigPath tests the XDG config path generation
func TestXDGConfigPath(t *testing.T) {
	// Test with XDG_CONFIG_HOME set
	originalXDGHome := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", originalXDGHome)
		os.Setenv("HOME", originalHome)
	}()

	// Test with XDG_CONFIG_HOME
	os.Setenv("XDG_CONFIG_HOME", "/custom/config")
	os.Setenv("HOME", "/home/user")
	path := getXDGConfigPath()
	assert.Equal(t, "/custom/config/krapp/config.yaml", path)

	// Test without XDG_CONFIG_HOME
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Setenv("HOME", "/home/user")
	path = getXDGConfigPath()
	assert.Equal(t, "/home/user/.config/krapp/config.yaml", path)
}

// TestMigrateLegacyConfig tests the legacy config migration
func TestMigrateLegacyConfig(t *testing.T) {
	tempDir := t.TempDir()
	legacyPath := filepath.Join(tempDir, ".krapp_config.yaml")
	newPath := filepath.Join(tempDir, ".config", "krapp", "config.yaml")

	// Create legacy config
	legacyConfig := Config{
		BaseDir: "/legacy/base",
		Editor:  "legacy-editor",
	}
	err := saveConfig(legacyPath, legacyConfig)
	assert.NoError(t, err)

	// Set paths to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Reset config paths to use new HOME
	ResetConfigPaths()

	// Set the config paths to use the temp directory
	SetConfigPaths(ConfigPaths{
		Global: newPath,
		Local:  filepath.Join(tempDir, ".krapp_config.yaml"),
	})

	// Run migration
	err = migrateLegacyConfig()
	assert.NoError(t, err)

	// Check that new config exists
	_, err = os.Stat(newPath)
	assert.NoError(t, err, "New config should exist")

	// Check that legacy config is removed
	_, err = os.Stat(legacyPath)
	assert.True(t, os.IsNotExist(err), "Legacy config should be removed")

	// Check that config content is preserved
	migratedConfig, err := loadConfig(newPath)
	assert.NoError(t, err)
	assert.Equal(t, legacyConfig.BaseDir, migratedConfig.BaseDir)
	assert.Equal(t, legacyConfig.Editor, migratedConfig.Editor)
}
