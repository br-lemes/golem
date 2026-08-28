package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/parser"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var equipCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "equip <name> <code>",
	Short: "Equip one or more items on a character",
	Long: `Equip one or more items on a character

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).EquipmentArgs(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		items, err := parser.Items(args[1:])
		if err != nil {
			return err
		}
		equipments, err := items.EquipSchemas()
		if err != nil {
			return err
		}
		_, err = routine.Equip(name, equipments)
		return err
	},
}

func init() {
	rootCmd.AddCommand(equipCmd)
}
