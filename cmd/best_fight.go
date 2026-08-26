package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type bestFightFlags struct {
	Name                    string `flag:"name" shorthand:"n" desc:"Use a real character as the base"`
	Level                   int    `flag:"level" default:"1" desc:"Fake character level (1-50)"`
	All                     bool   `flag:"all" desc:"Include all equippable equipment, even if not owned"`
	AllowDuplicateAdeptRing bool   `flag:"allow-duplicate-adept-ring" desc:"Allow using ring_of_the_adept twice"`
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
	_, ok := database.Monsters.Get(monster)
	if !ok {
		return fmt.Errorf("invalid monster: %s", monster)
	}
	if flags.Level < 1 || flags.Level > 50 {
		return fmt.Errorf("level must be between 1 and 50")
	}
	if flags.Name != "" && levelChanged {
		return fmt.Errorf("--name and --level are mutually exclusive")
	}
	return nil
}

func bestFightRun(monster string, flags bestFightFlags) error {
	monsterData, _ := database.Monsters.Get(monster)
	var result best.Result
	var err error
	if flags.Name != "" {
		result, err = best.FindFightByName(flags.Name, *monsterData, flags.All, flags.AllowDuplicateAdeptRing)
	} else {
		result, err = best.FindFightAtLevel(flags.Level, *monsterData, flags.All, flags.AllowDuplicateAdeptRing)
	}
	if err != nil {
		return err
	}
	return console.Auto(result)
}

func init() {
	bestCmd.AddCommand(bestFightCmd)
	err := utils.RegisterFlags[bestFightFlags](bestFightCmd)
	if err != nil {
		panic(err)
	}
}
