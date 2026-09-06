package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Access game database data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func databaseCommand[O, T any](use, short string, fetch func(O) (T, error)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := fetch(*new(O))
			if err != nil {
				return err
			}
			return console.Auto(result)
		},
	}
}

func init() {
	rootCmd.AddCommand(databaseCmd)
	commands := []*cobra.Command{
		databaseCommand("achievements", "Get All Achievements", api.Achievements),
		databaseCommand("effects", "Get All Effects", api.Effects),
		databaseCommand("events", "Get All Events", api.Events),
		databaseCommand("items", "Get All Items", api.Items),
		databaseCommand("maps", "Get All Maps", api.Maps),
		databaseCommand("monsters", "Get All Monsters", api.Monsters),
		databaseCommand("npcsDetails", "Get All Npcs", api.NpcsDetails),
		databaseCommand("resources", "Get All Resources", api.Resources),
		databaseCommand("tasksList", "Get All Tasks", api.TasksList),
	}
	databaseCmd.AddCommand(commands...)
}
