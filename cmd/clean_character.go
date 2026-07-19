package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var cleanCharacterCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(5),
	Use:   "character [names...]",
	Short: "Clean the character cache",
	Long: `Clean the character cache

Arguments:
  names   Names of the characters to clean.`,
	ValidArgsFunction: completion.CharacterName(5).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && console.Confirm("Clean all characters") {
			for _, character := range utils.GetCharacters() {
				cache.CleanCharacter(character)
			}
		}
		for _, character := range args {
			if !slices.Contains(utils.GetCharacters(), character) {
				return fmt.Errorf("character %s not found", character)
			}
		}
		for _, character := range args {
			cache.CleanCharacter(character)
		}
		return nil
	},
}

func init() {
	cleanCmd.AddCommand(cleanCharacterCmd)
}
