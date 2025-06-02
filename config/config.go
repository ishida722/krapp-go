package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseDir      string `yaml:"base_dir"`
	DailyNoteDir string `yaml:"daily_note_dir"`
	Inbox        string `yaml:"inbox_dir,omitempty"` // オプションのフィールド
}

var defaultConfig = Config{
	BaseDir:      "./notes",
	DailyNoteDir: "daily",
	Inbox:        "inbox", // デフォルトのInboxディレクトリ
}

var globalConfig *Config

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		// ファイルがなければデフォルト値を返す
		globalConfig = &defaultConfig
		return globalConfig, nil
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	globalConfig = &cfg
	return globalConfig, nil
}

func SaveConfig(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	return encoder.Encode(cfg)
}

func GetConfig() *Config {
	return globalConfig
}

// MarshalYAML はConfig構造体をYAMLバイト列に変換します。
func MarshalYAML(cfg *Config) ([]byte, error) {
	return yaml.Marshal(cfg)
}
