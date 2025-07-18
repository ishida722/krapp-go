package config

import (
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseDir              string         `yaml:"base_dir"`
	DailyNoteDir         string         `yaml:"daily_note_dir"`
	Inbox                string         `yaml:"inbox_dir"`
	Editor               string         `yaml:"editor"`
	WithAlwaysOpenEditor bool           `yaml:"with_always_open_editor"` // trueなら常にエディタを開く
	EditorOption         string         `yaml:"editor_option"`           // エディタのオプション
	DailyTemplate        map[string]any `yaml:"daily_template"`          // デイリーノート用テンプレート
	InboxTemplate        map[string]any `yaml:"inbox_template"`          // インボックスノート用テンプレート
}

var defaultConfig = Config{
	BaseDir:              "./notes",
	DailyNoteDir:         "daily",
	Inbox:                "inbox", // デフォルトのInboxディレクトリ
	Editor:               "vim",   // デフォルトのエディタ
	WithAlwaysOpenEditor: false,   // デフォルトでは常にエディタを開かない
	EditorOption:         "",      // デフォルトのエディタオプション
	DailyTemplate: map[string]any{
		"tags": []string{},
	},
	InboxTemplate: map[string]any{
		"tags":   []string{},
		"status": "new",
	},
}

type ConfigPaths struct {
	Global string // グローバル設定ファイルのパス
	Local  string // ローカル設定ファイルのパス
}

// getXDGConfigPath returns the XDG-compliant config path for global settings
func getXDGConfigPath() string {
	// Check XDG_CONFIG_HOME first
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "krapp", "config.yaml")
	}
	// Fall back to ~/.config/krapp/config.yaml
	return filepath.Join(os.Getenv("HOME"), ".config", "krapp", "config.yaml")
}

var defaultConfigPaths = ConfigPaths{
	Global: getXDGConfigPath(),
	Local:  ".krapp_config.yaml",
}

func GetDefaultConfigPaths() ConfigPaths {
	// デフォルトの設定ファイルパスを返す
	return defaultConfigPaths
}

var configPaths = GetDefaultConfigPaths()

func GetConfigPaths() (ConfigPaths, error) {
	return configPaths, nil
}

func SetConfigPaths(paths ConfigPaths) {
	configPaths = paths
}

func ResetConfigPaths() {
	// 設定ファイルのパスをデフォルトにリセット
	configPaths = GetDefaultConfigPaths()
}

// migrateLegacyConfig migrates old config file to new XDG-compliant location
func migrateLegacyConfig() error {
	// Check for legacy config file
	legacyPath := filepath.Join(os.Getenv("HOME"), ".krapp_config.yaml")
	if _, err := os.Stat(legacyPath); os.IsNotExist(err) {
		// Legacy config doesn't exist, nothing to migrate
		return nil
	}

	// Load legacy config
	legacyConfig, err := loadConfig(legacyPath)
	if err != nil {
		return err
	}

	// Create new config directory
	configDir := filepath.Dir(configPaths.Global)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Save config to new location
	if err := saveConfig(configPaths.Global, legacyConfig); err != nil {
		return err
	}

	// Remove legacy config file
	if err := os.Remove(legacyPath); err != nil {
		// Log warning but don't fail migration
		// TODO: Add proper logging when available
	}

	return nil
}

func makeHomeConfig() error {
	// First, try to migrate legacy config if it exists
	if err := migrateLegacyConfig(); err != nil {
		return err
	}

	// ファイルの情報を取得する,存在しない場合はエラーを返す
	_, err := os.Stat(configPaths.Global)
	// エラーが返ってこないので設定ファイルが存在する
	if err == nil {
		// ホームディレクトリに設定ファイルが存在する場合は何もしない
		return nil
	}
	// ファイルが存在するけど､エラーが発生した場合はそのエラーを返す
	if !os.IsNotExist(err) {
		return err // 存在しない以外のエラー
	}

	// 設定ファイルが存在しない場合、ディレクトリを作成
	configDir := filepath.Dir(configPaths.Global)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// ホームディレクトリに設定ファイルが存在しない場合はデフォルト設定を保存
	return saveConfig(configPaths.Global, GetDefaultConfig())
}

// expandHomePath expands ~ in path to the user's home directory
func expandHomePath(path string) string {
	if path == "" {
		return path
	}

	if path == "~" {
		return os.Getenv("HOME")
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.Getenv("HOME"), path[2:])
	}

	return path
}

func LoadConfig() (Config, error) {
	// 設定ファイルの存在確認と作成
	err := makeHomeConfig()
	if err != nil {
		// 設定ファイルの作成に失敗した場合はエラーを返す
		return Config{}, err
	}
	// グローバル設定の読み込み
	globalConfig, err := loadConfig(configPaths.Global)
	if err != nil {
		return Config{}, err
	}
	// ローカル設定の読み込み
	localConfig, err := loadConfig(configPaths.Local)
	if err != nil {
		// ローカル設定ファイルが存在しない場合はグローバル設定をそのまま返す
		// まずデフォルト設定とグローバル設定をマージ
		mergedConfig := MergeConfig(defaultConfig, globalConfig)
		// BaseDir内の~をホームディレクトリに展開
		mergedConfig.BaseDir = expandHomePath(mergedConfig.BaseDir)
		return mergedConfig, nil
	}

	// まずデフォルト設定とグローバル設定をマージ
	mergedGlobalConfig := MergeConfig(defaultConfig, globalConfig)

	// グローバル設定とローカル設定をマージ
	fixedConfig := MergeConfig(mergedGlobalConfig, localConfig)

	// BaseDir内の~をホームディレクトリに展開
	fixedConfig.BaseDir = expandHomePath(fixedConfig.BaseDir)

	return fixedConfig, nil
}

func loadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		// ファイルがなければnull値を返す
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	// yamlファイルをデコード
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		// デコードに失敗した場合はデフォルト値を返す
		return Config{}, err
	}
	return cfg, nil
}

func saveConfig(path string, cfg Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewEncoder(f).Encode(cfg)
}

// MarshalYAML はConfig構造体をYAMLバイト列に変換します。
func MarshalYAML(cfg *Config) ([]byte, error) {
	return yaml.Marshal(cfg)
}

func GetDefaultConfig() Config {
	// デフォルト設定を返す
	return defaultConfig
}

// MergeConfig はグローバル設定とローカル設定をマージします。
// ローカル設定の非ゼロ値でグローバル設定を上書きします。
// もしマージに失敗した場合はグローバル設定をそのまま返します。
// グローバル設定を基本として､ローカル設定で設定されている値で上書きします｡
// 例えば、ローカル設定でエディタが指定されていれば、グローバル設定のエディタを上書きします。
func MergeConfig(global, local Config) Config {
	// localの非ゼロ値でhomeを上書き
	if err := mergo.Merge(&local, global); err != nil {
		// マージに失敗した場合はグローバル設定をそのまま返す
		return global
	}
	return local
}
