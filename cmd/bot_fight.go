package cmd

import (
	"fmt"
	"time"

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

		character, err := apiCharacters(name)
		if err != nil {
			return err
		}

		fmt.Fprintf(writer, "[%s] Starting fight bot for %s\n",
			time.Now().Format("15:04:05"), name)

		for {
			character, err = handleHp(character)
			if err != nil {
				return err
			}

			character, err = handleInventory(character)
			if err != nil {
				return err
			}

			fmt.Fprintf(writer, "[%s] Current position: (%d, %d)\n",
				time.Now().Format("15:04:05"), character.X, character.Y)

			character, err = handleMap(character, code)
			if err != nil {
				return err
			}

			fmt.Fprintf(writer, "[%s] Starting fight...\n",
				time.Now().Format("15:04:05"))
			fightResult, err := apiActionFight(name, []string{})
			if err != nil {
				return fmt.Errorf("error during fight: %w", err)
			}
			character = fightResult.Characters[0]
			fight := fightResult.Fight
			fightStats := fight.Characters[0]
			var result string
			if fight.Result == "win" {
				result = "🏆 Fight won!"
			} else {
				result = "💀 Fight lost!"
			}
			fmt.Fprintln(writer, result)
			fmt.Fprintf(writer, "⚔️  XP gained: %d | HP remaining: %d\n",
				fightStats.Xp, character.Hp)

			if len(fightStats.Drops) > 0 {
				dropsStr := ""
				for _, drop := range fightStats.Drops {
					dropsStr += fmt.Sprintf("%dx %s, ", drop.Quantity, drop.Code)
				}
				fmt.Fprintf(writer, "🎁 Loot dropped: %s\n", dropsStr)
			}

			fmt.Fprintln(writer)
		}
	},
}

func init() {
	botCmd.AddCommand(botFightCmd)
}
