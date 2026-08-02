package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var sizesCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "sizes [route]",
	Short: "List GET route pagination size limits",
	Long: `List GET route pagination size limits

Arguments:
  route   Path of a specific route to list sizes for.`,
	ValidArgsFunction: completion.Custom(1, utils.GetRoutesCompletion).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			size, err := utils.GetSize(args[0])
			if err != nil {
				return err
			}
			return console.Auto(map[string]int{args[0]: size})
		}

		sizes, err := utils.GetSizes()
		if err != nil {
			return err
		}

		return console.Auto(sizes)
	},
}

func init() {
	rootCmd.AddCommand(sizesCmd)
}
