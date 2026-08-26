package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var bestGatheringCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "gathering <name> <code>",
	Short: "Find the best equipment for gathering",
	Long: `Find the best equipment for gathering

Arguments:
  name   Name of your character.
  code   The code of the resource.`,
	ValidArgsFunction: completion.CharacterName(1).Resource(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]
		flags, err := utils.ReadFlags[bestFlags](cmd)
		if err != nil {
			return err
		}
		err = bestGatheringValidate(code)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return bestGatheringRun(name, code, flags)
	},
}

func bestGatheringValidate(code string) error {
	_, ok := database.Resources.Get(code)
	if !ok {
		return fmt.Errorf("resource not found: %s", code)
	}
	return nil
}

func bestGatheringRun(name, code string, flags bestFlags) error {
	resource, _ := database.Resources.Get(code)
	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	skill := string(resource.Skill)
	priorities := []string{skill, "prospecting"}
	if gatheringSkillLevel(character, skill)-resource.Level <= 10 {
		priorities = []string{skill, "wisdom", "prospecting"}
	}
	items, err := best.FindEquipment(character, !flags.AllowDuplicateAdeptRing, priorities...)
	if err != nil {
		return err
	}
	return console.Auto(items)
}

func gatheringSkillLevel(c schemas.CharacterSchema, skill string) int {
	switch skill {
	case "alchemy":
		return c.AlchemyLevel
	case "fishing":
		return c.FishingLevel
	case "mining":
		return c.MiningLevel
	case "woodcutting":
		return c.WoodcuttingLevel
	}
	return 0
}

func init() {
	bestCmd.AddCommand(bestGatheringCmd)
	err := utils.RegisterFlags[bestFlags](bestGatheringCmd)
	if err != nil {
		panic(err)
	}
}
