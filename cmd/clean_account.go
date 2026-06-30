package cmd

import (
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/spf13/cobra"
)

var cleanAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Clean the account cache",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cache.CleanAccount()
		return nil
	},
}

func init() {
	cleanCmd.AddCommand(cleanAccountCmd)
}
