package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ishida722/krapp-go/config"
	"github.com/ishida722/krapp-go/usecase"

	"github.com/spf13/cobra"
)

type configAdapter struct{ *config.Config }

func (c *configAdapter) GetBaseDir() string      { return c.BaseDir }
func (c *configAdapter) GetDailyNoteDir() string { return c.DailyNoteDir }
func (c *configAdapter) GetInboxDir() string     { return c.Inbox }

func OpenFile(cmd *cobra.Command, config config.Config, filePath string) error {
	if config.WithAlwaysOpenEditor || cmd.Flags().Changed("edit") {
		err := usecase.OpenFile(config.Editor, filePath)
		if err != nil {
			return fmt.Errorf("ファイルを開く際にエラーが発生しました:%s", err)
		}
		return nil
	}
	return nil
}

func main() {
	// コンフィグのロード
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("設定ファイルの読み込みに失敗しました:", err)
		os.Exit(1)
	}

	adapter := &configAdapter{&cfg}

	var rootCmd = &cobra.Command{
		Use:     "krapp",
		Version: "0.2.1",
		Short:   "My awesome CLI tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("BaseDir: %s\nDailyNoteDir: %s\n", cfg.BaseDir, cfg.DailyNoteDir)
		},
	}

	var printConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Print current config as YAML",
		Run: func(cmd *cobra.Command, args []string) {
			// 設定内容をYAMLで出力
			yamlBytes, err := config.MarshalYAML(&cfg)
			if err != nil {
				fmt.Println("設定のYAML変換に失敗:", err)
				os.Exit(1)
			}
			fmt.Print(string(yamlBytes))
		},
	}
	rootCmd.AddCommand(printConfigCmd)

	var createDailyCmd = &cobra.Command{
		Use:     "create-daily",
		Short:   "Create today's daily note and print its path",
		Aliases: []string{"cd"},
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			filePath, err := usecase.CreateDailyNote(adapter, now)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(filePath)
			err = OpenFile(cmd, cfg, filePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	createDailyCmd.Flags().BoolP("edit", "e", false, "Open the note in editor after creation")
	rootCmd.AddCommand(createDailyCmd)

	var createInboxCmd = &cobra.Command{
		Use:     "create-inbox [title]",
		Short:   "Create a new inbox note with the given title and print its path",
		Aliases: []string{"ci"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := args[0]
			now := time.Now()
			filePath, err := usecase.CreateInboxNote(adapter, now, title)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(filePath)
			err = OpenFile(cmd, cfg, filePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	createInboxCmd.Flags().BoolP("edit", "e", false, "Open the note in editor after creation")
	rootCmd.AddCommand(createInboxCmd)

	var syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync notes",
		Run: func(cmd *cobra.Command, args []string) {
			usecase.SyncGit(cfg.BaseDir)
		},
	}
	rootCmd.AddCommand(syncCmd)

	var importNotes = &cobra.Command{
		Use:     "import-notes [directory]",
		Short:   "import notes from a directory",
		Aliases: []string{"in"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			directory := args[0]
			if err := usecase.ImportNotes(directory, filepath.Join(cfg.BaseDir, cfg.Inbox)); err != nil {
				fmt.Println("ノートのインポートに失敗しました:", err)
				os.Exit(1)
			} else {
				fmt.Println("ノートのインポートが完了しました")
			}
		},
	}
	rootCmd.AddCommand(importNotes)

	rootCmd.Execute()
}
