package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.Resources()
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbResourcesCmd)
}
