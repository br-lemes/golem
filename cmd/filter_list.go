package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var filterListCmd = &cobra.Command{
	Use:   "list [command]",
	Short: "List persistent output filters",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command := ""
		if len(args) == 1 {
			command = args[0]
		}
		return console.Auto(cache.ListOutputFilters(command))
	},
}

func init() {
	filterCmd.AddCommand(filterListCmd)
}
