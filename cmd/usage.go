package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/usage"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type usageFlags struct {
	Details bool `flag:"details" desc:"include monsters and effects"`
}

var usageCmd = &cobra.Command{
	Args:  cobra.ArbitraryArgs,
	Use:   "usage [code]...",
	Short: "Analyze equipment usage and shortages",
	Long: `Analyze equipment usage and shortages

Arguments:
  code   The code of the item.`,
	ValidArgsFunction: completion.Item(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[usageFlags](cmd)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return usageRun(args, flags)
	},
}

func usageRun(codes []string, flags usageFlags) error {
	result, err := usage.Evaluate(codes, flags.Details)
	if err != nil {
		return err
	}
	if !flags.Details {
		simple := make(map[string]bool, len(result))
		for code, evaluation := range result {
			simple[code] = evaluation.Needed
		}
		return console.Auto(simple)
	}
	return console.Auto(result)
}

func init() {
	rootCmd.AddCommand(usageCmd)
	err := utils.RegisterFlags[usageFlags](usageCmd)
	if err != nil {
		panic(err)
	}
}
