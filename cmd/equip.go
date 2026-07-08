package cmd

import (
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/task"
	"github.com/spf13/cobra"
)

var equipCmd = &cobra.Command{
	Use:   "equip <name> <code>",
	Short: "Equip Item",
	Args:  cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return config.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
		}
		return database.GetEquipments(), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := task.Equip(args[0], args[1:])
		return err
	},
}

func init() {
	rootCmd.AddCommand(equipCmd)
}
