package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbMapsCmd = &cobra.Command{
	Use:   "maps",
	Short: "Maps",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.Maps()
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbMapsCmd)
}
