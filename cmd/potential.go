package cmd

import (
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/potential"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type potentialFlags struct {
	CharacterLevel int      `flag:"character-level" desc:"Limit the simulated character level"`
	Details        bool     `flag:"details" desc:"Include monster and effect uses"`
	EquipmentLevel int      `flag:"equipment-level" desc:"Limit the available equipment level"`
	Loadouts       bool     `flag:"loadouts" desc:"Include combat loadouts"`
	Without        []string `flag:"without" desc:"Exclude items from the available equipment"`
}

var potentialCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "potential <code> [code...]",
	Short: "Analyze equipment's potential use",
	Long: `Analyze equipment's potential use

Arguments:
	code   The code of the item.`,
	ValidArgsFunction: completion.Custom(0, catalog.Items().Equipments().Keys).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[potentialFlags](cmd)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return potentialRun(args, flags)
	},
}

func potentialRun(codes []string, flags potentialFlags) error {
	result, err := potential.Evaluate(codes, flags.Without, flags.Loadouts, flags.CharacterLevel, flags.EquipmentLevel)
	if err != nil {
		return err
	}
	if !flags.Details && !flags.Loadouts {
		return console.Auto(potentialSimple(result))
	}
	return console.Auto(result)
}

type potentialSimpleResult struct {
	Used              bool     `json:"used"`
	SelectedEquipment []string `json:"selected_equipment,omitempty"`
}

func potentialSimple(result map[string]potential.Result) map[string]potentialSimpleResult {
	simple := make(map[string]potentialSimpleResult, len(result))
	for code, evaluation := range result {
		selected := make([]string, 0, len(evaluation.SelectedEquipment))
		for _, equipment := range evaluation.SelectedEquipment {
			selected = append(selected, equipment.Code)
		}
		simple[code] = potentialSimpleResult{
			Used:              evaluation.Used,
			SelectedEquipment: selected,
		}
	}
	return simple
}

func init() {
	rootCmd.AddCommand(potentialCmd)
	err := utils.RegisterFlags[potentialFlags](potentialCmd)
	if err != nil {
		panic(err)
	}
	err = potentialCmd.RegisterFlagCompletionFunc("without", completion.StringSlice(catalog.Items().Equipments().Keys))
	if err != nil {
		panic(err)
	}
}
