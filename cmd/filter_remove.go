package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var filterRemoveCmd = &cobra.Command{
	Use:   "remove <command> <kind> <pattern>",
	Short: "Remove a persistent output filter",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := validateFilterArgs(args, 3)
		if err != nil {
			return err
		}
		removed, err := cache.RemoveOutputFilter(args[0], args[1], args[2])
		if err != nil {
			return err
		}
		if !removed {
			console.Debugf("output filter not found: command=%s kind=%s pattern=%s\n", args[0], args[1], args[2])
		}
		return nil
	},
}

func init() {
	filterCmd.AddCommand(filterRemoveCmd)
}
