package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var gatheringCmd = &cobra.Command{
	Use:   "gathering <name> <code>",
	Short: "Gather resources continuously",
	Long: `Gather resources continuously

Arguments:
  name   Name of your character.
  code   The code of the resource.`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			characters := config.GetCharacters()
			return characters, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			codes := database.GetResourceCodes()
			return codes, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		resource, found := database.GetResource(code)
		if !found {
			return fmt.Errorf("resource %s not found", code)
		}

		if !config.ConfirmSkill(name, string(resource.Skill)) {
			cmd.SilenceUsage = true
			return fmt.Errorf("operation cancelled")
		}

		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		routine.Cooldown(character)

		for {
			character, err = routine.Inventory(character, []string{})
			if err != nil {
				return err
			}

			character, err = routine.Move(character, code)
			if err != nil {
				return err
			}

			skill, err := api.MyActionGathering(name)
			if err != nil {
				return err
			}

			character = skill.Character
		}
	},
}

func init() {
	rootCmd.AddCommand(gatheringCmd)
}
