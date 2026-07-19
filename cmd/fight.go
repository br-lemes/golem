package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var fightCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "fight <name> <code>",
	Short: "Fight continuously",
	Long: `Fight continuously

Arguments:
  name   Name of your character.
  code   The code of the monster.`,
	ValidArgsFunction: completion.CharacterName(1).Monster(1).Build(),
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
		if taskCompleted(character, code) {
			return nil
		}
		if character.Level < monster.Level {
			console.Printf("Your level %d < monster level %d\n",
				character.Level, monster.Level)
			if !console.Confirm("Do you want to continue?") {
				cmd.SilenceUsage = true
				return fmt.Errorf("operation cancelled")
			}
		}

		for {
			prepare(character, monster)

			fightResult, err := api.MyActionFight(name, []string{})
			if err != nil {
				return fmt.Errorf("error during fight: %w", err)
			}
			character = fightResult.Characters[0]

			if taskCompleted(character, code) {
				return nil
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fightCmd)
}

func prepare(character schemas.CharacterSchema, monster schemas.MonsterSchema) error {
	routine.Cooldown(character)
	minHp := monster.Hp + (monster.Hp * 20 / 100)
	if minHp > character.MaxHp {
		minHp = character.MaxHp
	}
	character, err := routine.Hp(character, minHp)
	if err != nil {
		return err
	}
	character, err = routine.Inventory(character, []string{"food"})
	if err != nil {
		return err
	}
	_, err = routine.Move(character, monster.Code)
	if err != nil {
		return err
	}
	return nil
}

func taskCompleted(character schemas.CharacterSchema, code string) bool {
	if character.TaskType == "monsters" &&
		character.Task == code &&
		character.TaskProgress == character.TaskTotal {
		console.Printf("  Task completed\n")
		return true
	}
	return false
}
