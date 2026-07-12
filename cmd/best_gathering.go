package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestGatheringCmd = &cobra.Command{
	Use:   "gathering <name> <skill>",
	Short: "Find the best equipment for gathering",
	Long: `Find the best equipment for gathering

Arguments:
  name   Name of your character.
  skill  The gathering skill (alchemy, fishing, mining, woodcutting).
`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			validSkills := database.GetEnum("GatheringSkill")
			return validSkills, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		skill := args[1]
		validSkills := database.GetEnum("GatheringSkill")
		if !slices.Contains(validSkills, skill) {
			return fmt.Errorf("invalid skill: %s", skill)
		}
		character, err := api.Characters(name)
		if err != nil {
			return err
		}

		items, err := utils.BestFinder(character, skill, map[string]int{
			skill: 16, "prospecting": 8, "wisdom": 4, "inventory_space": 2})
		if err != nil {
			return err
		}
		console.Auto(items)
		return nil
	},
}

func init() {
	bestCmd.AddCommand(bestGatheringCmd)
}
