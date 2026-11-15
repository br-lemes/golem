package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type bestFightFlags struct {
	Name                    string   `flag:"name" shorthand:"n" desc:"Use a real character as the base"`
	Level                   int      `flag:"level" default:"1" desc:"Fake character level (1-50)"`
	All                     bool     `flag:"all" desc:"Include all equippable equipment, even if not owned"`
	Full                    bool     `flag:"full" desc:"Show the complete final equipment"`
	AllowDuplicateAdeptRing bool     `flag:"allow-duplicate-adept-ring" desc:"Allow using ring_of_the_adept twice"`
	Without                 []string `flag:"without" desc:"Forbid items from the selection, even if owned"`
	With                    []string `flag:"with" desc:"Include items in the selection, even if not owned"`
	Use                     []string `flag:"use" desc:"Force using specific items in the selection"`
}

var bestFightCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "fight <monster>",
	Short: "Find equipment using the current combat simulator",
	Long: `Find equipment using the current combat simulator

Arguments:
  monster   The code of the monster.`,
	ValidArgsFunction: completion.Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		monster := args[0]
		flags, err := utils.ReadFlags[bestFightFlags](cmd)
		if err != nil {
			return err
		}
		err = bestFightValidate(monster, flags, cmd.Flags().Changed("level"))
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return bestFightRun(monster, flags)
	},
}

func bestFightValidate(monster string, flags bestFightFlags, levelChanged bool) error {
	_, ok := catalog.Monsters.Get(monster)
	if !ok {
		return fmt.Errorf("invalid monster: %s", monster)
	}
	err := validateGameLevel(flags.Level)
	if err != nil {
		return err
	}
	if flags.Name != "" && levelChanged {
		return fmt.Errorf("--name and --level are mutually exclusive")
	}
	return bestFightValidateItems(flags)
}

func bestFightValidateItems(flags bestFightFlags) error {
	for _, code := range flags.Without {
		_, exists := catalog.Items().Get(code)
		if !exists {
			return fmt.Errorf("invalid item: %s", code)
		}
	}
	keep := make(map[string]bool)
	for _, code := range flags.With {
		_, exists := catalog.Items().Get(code)
		if !exists {
			return fmt.Errorf("invalid item: %s", code)
		}
		keep[code] = true
	}
	for _, code := range flags.Use {
		_, exists := catalog.Items().Get(code)
		if !exists {
			return fmt.Errorf("invalid item: %s", code)
		}
		keep[code] = true
	}
	for _, code := range flags.Without {
		if keep[code] {
			return fmt.Errorf("item cannot be forbidden and forced: %s", code)
		}
	}
	return nil
}

func bestFightRun(monster string, flags bestFightFlags) error {
	monsterData, _ := catalog.Monsters.Get(monster)
	options := best.FightOptions{
		All:       flags.All,
		AdeptRing: flags.AllowDuplicateAdeptRing,
		Without:   flags.Without,
		With:      flags.With,
		Use:       flags.Use,
	}
	var result best.Result
	var err error
	if flags.Name != "" {
		result, err = best.FindFightByName(flags.Name, *monsterData, options)
	} else {
		result, err = best.FindFightAtLevel(flags.Level, *monsterData, options)
	}
	if err != nil {
		return err
	}
	if flags.Full {
		result = bestFightFullResult(result)
	}
	return console.Auto(result)
}

func bestFightFullResult(result best.Result) best.Result {
	equipment := make(map[string]string)
	utilities := make(map[string]string)
	for slot, code := range result.FinalEquipment {
		switch slot {
		case "utility1", "utility2":
			utilities[slot] = code
		default:
			equipment[slot] = code
		}
	}
	result.Equipment = equipment
	result.Utilities = utilities
	return result
}

func init() {
	bestCmd.AddCommand(bestFightCmd)
	err := utils.RegisterFlags[bestFightFlags](bestFightCmd)
	if err != nil {
		panic(err)
	}
	err = bestFightCmd.RegisterFlagCompletionFunc("name", completion.SingleValue(cache.GetCharacters))
	if err != nil {
		panic(err)
	}
	for _, name := range []string{"without", "with", "use"} {
		err = bestFightCmd.RegisterFlagCompletionFunc(name, completion.UniqueList(catalog.Items().Equipments().Keys))
		if err != nil {
			panic(err)
		}
	}
}
