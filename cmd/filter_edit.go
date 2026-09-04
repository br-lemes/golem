package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/spf13/cobra"
)

var filterEditCmd = &cobra.Command{
	Use:   "edit <command> <kind> <old-pattern> <new-pattern>",
	Short: "Edit a persistent output filter",
	Args:  cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateFilterArgs(args[:3], 3)
		if err != nil {
			return err
		}
		if args[3] == "" {
			return fmt.Errorf("new pattern cannot be empty")
		}
		return cache.EditOutputFilter(args[0], args[1], args[2], args[3])
	},
}

func init() {
	filterCmd.AddCommand(filterEditCmd)
}
