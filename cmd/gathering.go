package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/spf13/cobra"
)

var gatheringCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "gathering <name> <code>",
	Short: "Gather resources continuously",
	Long: `Gather resources continuously

Arguments:
  name   Name of your character.
  code   The code of the resource.`,
	ValidArgsFunction: completion.CharacterName(1).Resource(1).Build(),
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
