package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbNpcsDetails = &cobra.Command{
	Use:   "npcsDetails",
	Short: "Npcs Details",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.NpcsDetails()
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbNpcsDetails)
}
