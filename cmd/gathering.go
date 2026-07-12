package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/utils"
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
			return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
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
