package krapp

import (
	"fmt"
	"os"

	"github.com/ishida722/krapp-go/usecase"
	"github.com/spf13/cobra"
)

func importIssuesCmd() *cobra.Command {
	var (
		repo    string
		dryRun  bool
		noClose bool
	)

	cmd := &cobra.Command{
		Use:     "import-issues",
		Short:   "Import GitHub issues as inbox notes",
		Aliases: []string{"ii"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfigAdapter()
			client := &usecase.GHClient{}

			options := usecase.ImportOptions{
				Repo:    repo,
				DryRun:  dryRun,
				NoClose: noClose,
			}

			if err := usecase.ImportGitHubIssues(cfg, client, options); err != nil {
				fmt.Printf("GitHub issueのインポートに失敗しました: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("GitHub issueのインポートが完了しました")
		},
	}

	cmd.Flags().StringVar(&repo, "repo", "", "Repository (owner/name)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Don't actually close issues")
	cmd.Flags().BoolVar(&noClose, "no-close", false, "Import issues without closing them")

	return cmd
}