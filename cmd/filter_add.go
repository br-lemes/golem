package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/spf13/cobra"
)

var filterAddCmd = &cobra.Command{
	Use:   "add <command> <kind> <pattern>",
	Short: "Add a persistent output filter",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateFilterArgs(args, 3)
		if err != nil {
			return err
		}
		return cache.AddOutputFilter(args[0], args[1], args[2])
	},
}

func init() {
	filterCmd.AddCommand(filterAddCmd)
}
