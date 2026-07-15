package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestFightCmd = &cobra.Command{
	Use:   "fight <name> <monster>",
	Short: "Find the best equipment for fighting",
	Long: `Find the best equipment for fighting

Arguments:
  name     Name of your character.
  monster  Code of the monster to fight.
`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			return database.GetMonsterCodes(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},

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
		items, err := utils.BestFight(character, monster)
		if err != nil {
			return err
		}
		console.Auto(items)
		return nil
	},
}

func init() {
	bestCmd.AddCommand(bestFightCmd)
}
