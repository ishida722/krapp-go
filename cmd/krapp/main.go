package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ishida722/krapp-go/config"

	"github.com/spf13/cobra"
)

func main() {
	// 1. カレントディレクトリの設定ファイルを優先
	localConfigPath := ".krapp_config.yaml"
	homeConfigPath := filepath.Join(os.Getenv("HOME"), ".krapp_config.yaml")
	configPath := homeConfigPath
	if _, err := os.Stat(localConfigPath); err == nil {
		configPath = localConfigPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("設定ファイルの読み込みに失敗しました:", err)
		os.Exit(1)
	}
	var rootCmd = &cobra.Command{
		Use:   "krapp",
		Short: "My awesome CLI tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("BaseDir: %s\nDailyNoteDir: %s\n", cfg.BaseDir, cfg.DailyNoteDir)
		},
	}

	var printConfigCmd = &cobra.Command{
		Use:   "print-config",
		Short: "Print current config as YAML",
		Run: func(cmd *cobra.Command, args []string) {
			// 設定内容をYAMLで出力
			yamlBytes, err := config.MarshalYAML(cfg)
			if err != nil {
				fmt.Println("設定のYAML変換に失敗:", err)
				os.Exit(1)
			}
			fmt.Print(string(yamlBytes))
		},
	}
	rootCmd.AddCommand(printConfigCmd)

	var createDailyCmd = &cobra.Command{
		Use:   "create-daily",
		Short: "Create today's daily note and print its path",
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			year := now.Format("2006")
			month := now.Format("01")
			date := now.Format("2006-01-02")
			dir := filepath.Join(cfg.BaseDir, cfg.DailyNoteDir, year, month)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Println("ディレクトリ作成に失敗:", err)
				os.Exit(1)
			}
			filePath := filepath.Join(dir, date+".md")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				f, err := os.Create(filePath)
				if err != nil {
					fmt.Println("ファイル作成に失敗:", err)
					os.Exit(1)
				}
				f.Close()
			}
			fmt.Println(filePath)
		},
	}
	rootCmd.AddCommand(createDailyCmd)

	var createInboxCmd = &cobra.Command{
		Use:   "create-inbox [title]",
		Short: "Create a new inbox note with the given title and print its path",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := args[0]
			now := time.Now()
			date := now.Format("2006-01-02")
			filename := fmt.Sprintf("%s-%s.md", date, title)
			dir := filepath.Join(cfg.BaseDir, cfg.Inbox)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Println("Inboxディレクトリ作成に失敗:", err)
				os.Exit(1)
			}
			filePath := filepath.Join(dir, filename)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				f, err := os.Create(filePath)
				if err != nil {
					fmt.Println("Inboxノート作成に失敗:", err)
					os.Exit(1)
				}
				f.Close()
			}
			fmt.Println(filePath)
		},
	}
	rootCmd.AddCommand(createInboxCmd)

	rootCmd.Execute()
}
