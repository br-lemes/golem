package cmd

import (
	_ "embed"
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [name]",
	Short: "Test commands",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		item, exists := database.GetItem(name)
		if exists {
			fmt.Printf("Item: %s (%s)\n", name, item.Type)
		}
		resource, exists := database.GetResource(name)
		if exists {
			fmt.Printf("Resource: %s (%s)\n", name, resource.Skill)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
