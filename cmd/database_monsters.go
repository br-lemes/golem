package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbMonstersCmd = &cobra.Command{
	Use:   "monsters",
	Short: "Monsters",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.Monsters()
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbMonstersCmd)
}
