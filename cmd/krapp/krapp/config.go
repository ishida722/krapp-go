package krapp

import (
	"fmt"
	"os"

	"github.com/ishida722/krapp-go/config"
	"github.com/spf13/cobra"
)

func printConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Print current config as YAML",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig()
			yamlBytes, err := config.MarshalYAML(&cfg)
			if err != nil {
				fmt.Println("設定のYAML変換に失敗:", err)
				os.Exit(1)
			}
			fmt.Print(string(yamlBytes))
		},
	}
}
