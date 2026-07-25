package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "simulate <name> <monster>",
	Short: "Simulate a fight against a monster",
	Long: `Simulate a fight against a monster

Arguments:
  name     Name of your character.
  monster  Code of the monster to fight.`,
	ValidArgsFunction: completion.CharacterName(1).Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		monsterCode := args[1]
		monster, exists := database.GetMonster(monsterCode)
		if !exists {
			return fmt.Errorf("invalid monster: %s", monsterCode)
		}
		character, err := api.Characters(name)
		if err != nil {
			return err
		}

		overrides := map[string]string{}
		for _, slot := range utils.FightSlots() {
			if !cmd.Flags().Changed(slot) {
				continue
			}
			v, err := cmd.Flags().GetString(slot)
			if err != nil {
				return err
			}
			overrides[slot] = v
		}

		report, err := utils.FightSimulate(character, monster, overrides)
		if err != nil {
			return err
		}
		console.Auto(report)
		return nil
	},
}

func init() {
	for _, slot := range utils.FightSlots() {
		usage := fmt.Sprintf("Override the %s slot for this simulation", slot)
		simulateCmd.Flags().String(slot, "", usage)
		simulateCmd.RegisterFlagCompletionFunc(slot, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return database.GetItemCodes(), cobra.ShellCompDirectiveNoFileComp
		})
	}
	rootCmd.AddCommand(simulateCmd)
}
