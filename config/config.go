package config

import (
	"os"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseDir      string `yaml:"base_dir"`
	DailyNoteDir string `yaml:"daily_note_dir"`
	Inbox        string `yaml:"inbox_dir"` // オプションのフィールド
	Editor       string `yaml:"editor"`    // オプションのフィールド
}

var defaultConfig = Config{
	BaseDir:      "./notes",
	DailyNoteDir: "daily",
	Inbox:        "inbox", // デフォルトのInboxディレクトリ
	Editor:       "vim",   // デフォルトのエディタ
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		// ファイルがなければデフォルト値を返す
		return defaultConfig, nil
	}
	defer f.Close()

	var cfg Config
	// yamlファイルをデコード
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		// デコードに失敗した場合はデフォルト値を返す
		return defaultConfig, err
	}
	return cfg, nil
}

func SaveConfig(path string, cfg Config) error {
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

func MergeConfig(local, home Config) Config {
	// localの非ゼロ値でhomeを上書き
	_ = mergo.Merge(home, local, mergo.WithOverride)
	return home
}
