package usecase

func load_config(){
	// カレントディレクトリの設定ファイル
	localConfigPath := ".krapp_config.yaml"
	// ホームディレクトリの設定ファイル
	homeConfigPath := filepath.Join(os.Getenv("HOME"), ".krapp_config.yaml")
	configPath := homeConfigPath
}