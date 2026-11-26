package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type bestFlags struct {
	All                     bool     `flag:"all" desc:"Include all equippable equipment, even if not owned"`
	AllowDuplicateAdeptRing bool     `flag:"allow-duplicate-adept-ring" desc:"Allow recommending ring_of_the_adept twice for this character"`
	Full                    bool     `flag:"full" desc:"Show the complete final equipment"`
	Loadout                 bool     `flag:"loadout" desc:"Output the selected loadout in simulation file format instead of the result"`
	With                    []string `flag:"with" desc:"Include items in the selection, even if not owned"`
	Without                 []string `flag:"without" desc:"Forbid items from the selection, even if owned"`
}

var bestCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(2),
	Use:   "best <name> <effect> [effect...]",
	Short: "Find the best equipment for the specified effects",
	Long: `Find the best equipment for the specified effects

Arguments:
  name     Name of your character.
  effect   The code of the effect.`,
	ValidArgsFunction: completion.CharacterName(1).Custom(0, catalog.Effects().Equipments().Keys).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		flags, err := utils.ReadFlags[bestFlags](cmd)
		if err != nil {
			return err
		}
		err = bestValidateItems(flags)
		if err != nil {
			return err
		}
		effects, err := best.NormalizePriorities(args[1:])
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return bestRun(name, flags, effects)
	},
}

func bestRun(name string, flags bestFlags, effects []string) error {
	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	items, err := best.FindEquipment(character, best.EquipmentOptions{
		UniqueAdeptRing: !flags.AllowDuplicateAdeptRing,
		Priorities:      effects,
		All:             flags.All,
		Without:         flags.Without,
		With:            flags.With,
		Full:            flags.Full || flags.Loadout,
	})
	if err != nil {
		return fmt.Errorf("find best equipment: %w", err)
	}
	return bestPrint(character, flags, items)
}

func bestValidateItems(flags bestFlags) error {
	forbidden := make(map[string]bool)
	for _, code := range flags.Without {
		_, exists := catalog.Items().Get(code)
		if !exists {
			return fmt.Errorf("invalid item: %s", code)
		}
		forbidden[code] = true
	}
	for _, code := range flags.With {
		item, exists := catalog.Items().Get(code)
		if !exists {
			return fmt.Errorf("invalid item: %s", code)
		}
		_, equippable := catalog.EquipmentTypeToSlots[item.Type]
		if !equippable {
			return fmt.Errorf("item cannot be equipped: %s", code)
		}
		if forbidden[code] {
			return fmt.Errorf("item cannot be forbidden and included: %s", code)
		}
	}
	return nil
}

func bestPrint(character schemas.CharacterSchema, flags bestFlags, items map[string]best.BestResult) error {
	if flags.Loadout {
		return console.Auto(bestPriorityLoadout(character, items))
	}
	return console.Auto(items)
}

func bestPriorityLoadout(character schemas.CharacterSchema, items map[string]best.BestResult) []schemas.FakeCharacterSchema {
	loadout := schemas.FakeCharacterSchema{Level: character.Level}
	for slot, item := range items {
		if item.Code == "" {
			continue
		}
		fight.SetSlot(&loadout, slot, item.Code)
		if slot == "utility1" || slot == "utility2" {
			fight.SetQuantity(&loadout, slot, bestFightUtilityQuantity)
		}
	}
	return []schemas.FakeCharacterSchema{loadout}
}

func registerBestFlags(cmd *cobra.Command) {
	err := utils.RegisterFlags[bestFlags](cmd)
	if err != nil {
		panic(err)
	}
	for _, name := range []string{"without", "with"} {
		err = cmd.RegisterFlagCompletionFunc(name, completion.UniqueList(catalog.Items().Equipments().Keys))
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(bestCmd)
	registerBestFlags(bestCmd)
}
