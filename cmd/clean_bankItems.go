package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/spf13/cobra"
)

var cleanBankItemsCmd = &cobra.Command{
	Use:   "bankItems",
	Short: "Clean the bank items cache",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cache.CleanBankItems()
		return nil
	},
}

func init() {
	cleanCmd.AddCommand(cleanBankItemsCmd)
}
