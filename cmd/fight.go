package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var fightCmd = &cobra.Command{
	Use:   "fight <name> <code>",
	Short: "Fight continuously",
	Long: `Fight continuously

Arguments:
  name   Name of your character.
  code   The code of the monster.`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			characters := config.GetCharacters()
			return characters, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			codes := database.GetMonsterCodes()
			return codes, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		monster, found := database.GetMonster(code)
		if !found {
			return fmt.Errorf("monster %s not found", code)
		}

		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		routine.Cooldown(character)

		if character.Level < monster.Level {
			console.Printf("Your level %d < monster level %d\n",
				character.Level, monster.Level)
			if !console.Confirm("Do you want to continue?") {
				cmd.SilenceUsage = true
				return fmt.Errorf("operation cancelled")
			}
		}

		minHp := monster.Hp + (monster.Hp * 20 / 100)
		if minHp > character.MaxHp {
			minHp = character.MaxHp
		}
		for {
			character, err = routine.Hp(character, minHp)
			if err != nil {
				return err
			}

			character, err = routine.Inventory(character, []string{"food"})
			if err != nil {
				return err
			}

			character, err = routine.Move(character, code)
			if err != nil {
				return err
			}

			fightResult, err := api.MyActionFight(name, []string{})
			if err != nil {
				return fmt.Errorf("error during fight: %w", err)
			}
			character = fightResult.Characters[0]
		}
	},
}

func init() {
	rootCmd.AddCommand(fightCmd)
}
