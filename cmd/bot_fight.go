package cmd

import (
	"fmt"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/task"
	"github.com/spf13/cobra"
)

var botFightCmd = &cobra.Command{
	Use:   "fight <name> <code>",
	Short: "Fight continuously",
	Long: `Fight continuously

Arguments:
  name   Name of your character.
  code   The code of the monster.`,
	Args: cobra.ExactArgs(2),
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
		task.Cooldown(character)

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
			character, err = task.Hp(character, minHp)
			if err != nil {
				return err
			}

			character, err = task.Inventory(character, true)
			if err != nil {
				return err
			}

			character, err = task.Move(character, code)
			if err != nil {
				return err
			}

			fightResult, err := api.MyActionFight(name, []string{})
			if err != nil {
				return fmt.Errorf("error during fight: %w", err)
			}
			character = fightResult.Characters[0]

			drops := []string{}
			for _, item := range fightResult.Fight.Characters[0].Drops {
				dropStr := fmt.Sprintf("%dx %s", item.Quantity, item.Code)
				drops = append(drops, dropStr)
			}

			if fightResult.Fight.Result == "win" {
				console.Printf("  XP gained: %d",
					fightResult.Fight.Characters[0].Xp)
				if len(drops) > 0 {
					dropsStr := strings.Join(drops, ", ")
					console.Printf(", Drops: %s", dropsStr)
				}
				console.Printf("\n")
			} else {
				console.Printf("  💀 Fight lost!\n")
			}
		}
	},
}

func init() {
	botCmd.AddCommand(botFightCmd)
}
