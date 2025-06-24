package krapp

import (
	"fmt"
	"os"
	"time"

	"github.com/ishida722/krapp-go/usecase"
	"github.com/spf13/cobra"
)

func createInboxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-inbox [title]",
		Short:   "Create a new inbox note with the given title and print its path",
		Aliases: []string{"ci"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig()
			adapter := &configAdapter{&cfg}

			title := args[0]
			now := time.Now()
			filePath, err := usecase.CreateInboxNote(adapter, now, title)
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
