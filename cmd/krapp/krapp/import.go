package krapp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ishida722/krapp-go/usecase"
	"github.com/spf13/cobra"
)

func importCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "import-notes [directory]",
		Short:   "import notes from a directory",
		Aliases: []string{"in"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig()
			directory := args[0]
			if err := usecase.ImportNotes(directory, filepath.Join(cfg.BaseDir, cfg.Inbox)); err != nil {
				fmt.Println("ノートのインポートに失敗しました:", err)
				os.Exit(1)
			} else {
				fmt.Println("ノートのインポートが完了しました")
			}
		},
	}
}
