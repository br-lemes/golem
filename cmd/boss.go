package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var bossCmd = &cobra.Command{
	Args:  cobra.RangeArgs(3, 4),
	Use:   "boss <name> <code> <participants>...",
	Short: "Boss fight continuously",
	Long: `Boss fight continuously

Arguments:
  name           Name of your character.
  code           The code of the boss.
  participants   Names of other characters in your account to fight alongside.`,
	ValidArgsFunction: completion.CharacterName(1).Boss(1).CharacterName(2).
		Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]
		participants := args[2:]
		switch len(participants) {
		case 1:
			if participants[0] == name {
				return fmt.Errorf("participants cannot include the main character")
			}
		case 2:
			if participants[0] == participants[1] {
				return fmt.Errorf("participants cannot include the same character twice")
			}
			if participants[0] == name || participants[1] == name {
				return fmt.Errorf("participants cannot include the main character")
			}
		}
		boss, found := database.GetMonster(code)
		if !found {
			return fmt.Errorf("boss %s not found", code)
		}
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}
		totalCharacters := len(participants) + 1
		charMap := make(map[string]schemas.CharacterSchema, totalCharacters)
		for _, character := range characters {
			charMap[character.Name] = character
		}
		_, found = charMap[name]
		if !found {
			return fmt.Errorf("character %s not found", name)
		}
		for _, participant := range participants {
			_, found := charMap[participant]
			if !found {
				return fmt.Errorf("character %s not found", participant)
			}
		}
		minLevel := charMap[name].Level
		for _, participant := range participants {
			if charMap[participant].Level < minLevel {
				minLevel = charMap[participant].Level
			}
		}
		if minLevel < boss.Level {
			console.Printf("Your minimum level %d < boss level %d\n",
				minLevel, boss.Level)
			if !console.Confirm("Do you want to continue?") {
				cmd.SilenceUsage = true
				return fmt.Errorf("operation cancelled")
			}
		}

		for {
			var g errgroup.Group
			g.Go(func() error { return prepare(charMap[name], boss) })
			for _, p := range participants {
				g.Go(func() error { return prepare(charMap[p], boss) })
			}
			err := g.Wait()
			if err != nil {
				return err
			}

			fightResult, err := api.MyActionFight(name, participants)
			if err != nil {
				return fmt.Errorf("error during fight: %w", err)
			}
			for _, character := range fightResult.Characters {
				charMap[character.Name] = character
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(bossCmd)
}
