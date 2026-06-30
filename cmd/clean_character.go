package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var cleanCharacterCmd = &cobra.Command{
	Use:   "character",
	Short: "Clean the character cache",
	Args:  cobra.MaximumNArgs(5),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && console.Confirm("Clean all characters") {
			for _, character := range config.GetCharacters() {
				cache.CleanCharacter(character)
			}
		}
		for _, character := range args {
			_, found := config.Characters[character]
			if !found {
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
