package cmd

import (
	_ "embed"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		item, exists := database.GetItem("ash_wood")
		if !exists {
			return nil
		}
		println(item.Type)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
