package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var cleanCharacterCmd = &cobra.Command{
	Use:   "character",
	Short: "Clean the character cache",
	Args:  cobra.MaximumNArgs(5),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
	},
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
