package krapp

import (
	"fmt"
	"os"

	"github.com/ishida722/krapp-go/config"
	"github.com/spf13/cobra"
)

var cfg config.Config

type configAdapter struct{ *config.Config }

func (c *configAdapter) GetBaseDir() string      { return c.BaseDir }
func (c *configAdapter) GetDailyNoteDir() string { return c.DailyNoteDir }
func (c *configAdapter) GetInboxDir() string     { return c.Inbox }

var rootCmd = &cobra.Command{
	Use:     "krapp",
	Version: "0.2.1",
	Short:   "My awesome CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("BaseDir: %s\nDailyNoteDir: %s\n", cfg.BaseDir, cfg.DailyNoteDir)
	},
}

func Execute() error {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Println("設定ファイルの読み込みに失敗しました:", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(printConfigCmd())
	rootCmd.AddCommand(createDailyCmd())
	rootCmd.AddCommand(createInboxCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(importCmd())

	return rootCmd.Execute()
}

func getConfig() config.Config {
	return cfg
}
