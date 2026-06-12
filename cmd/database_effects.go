package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbEffectsCmd = &cobra.Command{
	Use:   "effects",
	Short: "Effects",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.Effects()
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbEffectsCmd)
}
