package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbItemsCmd = &cobra.Command{
	Use:   "items",
	Short: "Items",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.Items()
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbItemsCmd)
}
