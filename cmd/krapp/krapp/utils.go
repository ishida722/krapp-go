package krapp

import (
	"fmt"

	"github.com/ishida722/krapp-go/config"
	"github.com/ishida722/krapp-go/usecase"
	"github.com/spf13/cobra"
)

func openFile(cmd *cobra.Command, config config.Config, filePath string) error {
	if config.WithAlwaysOpenEditor || cmd.Flags().Changed("edit") {
		err := usecase.OpenFile(config.Editor, filePath, config.EditorOption)
		if err != nil {
			return fmt.Errorf("ファイルを開く際にエラーが発生しました:%s", err)
		}
		return nil
	}
	return nil
}
