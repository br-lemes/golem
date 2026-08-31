package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/utils"
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
		resource, found := database.Resources.Get(code)
		if !found {
			return fmt.Errorf("resource %s not found", code)
		}

		character, err := api.Characters(name)
		if err != nil {
			return err
		}

		skillLevel, _ := utils.GetCharacterGatheringSkillLevel(character, string(resource.Skill))
		if skillLevel < resource.Level {
			return fmt.Errorf("character %s level too low. Required: %d, Current %d", name, resource.Level, skillLevel)
		}
		routine.Cooldown(character)

		equipments, err := best.FindEquipmentSchemas(character, best.EquipmentOptions{
			UniqueAdeptRing: true,
			Priorities:      best.GatheringPriorities(character, resource),
		})
		if err != nil {
			return err
		}
		character, err = routine.Equip(name, equipments)
		if err != nil {
			return err
		}

		for {
			character, err = routine.Inventory(character, []string{})
			if err != nil {
				return err
			}

			_, err = routine.Move(character, code)
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
