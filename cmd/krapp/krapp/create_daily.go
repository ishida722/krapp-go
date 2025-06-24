package krapp

import (
	"fmt"
	"os"
	"time"

	"github.com/ishida722/krapp-go/usecase"
	"github.com/spf13/cobra"
)

func createDailyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-daily",
		Short:   "Create today's daily note and print its path",
		Aliases: []string{"cd"},
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig()
			adapter := &configAdapter{&cfg}

			now := time.Now()
			filePath, err := usecase.CreateDailyNote(adapter, now)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(filePath)

			err = openFile(cmd, cfg, filePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().BoolP("edit", "e", false, "Open the note in editor after creation")
	return cmd
}
