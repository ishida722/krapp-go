package krapp

import (
	"github.com/ishida722/krapp-go/usecase"
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync notes",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig()
			usecase.SyncGit(cfg.BaseDir)
		},
	}
}
