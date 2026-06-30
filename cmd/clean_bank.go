package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/spf13/cobra"
)

var cleanBankCmd = &cobra.Command{
	Use:   "bank",
	Short: "Clean the bank cache",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cache.CleanBank()
		return nil
	},
}

func init() {
	cleanCmd.AddCommand(cleanBankCmd)
}
