package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type refreshFlags struct {
	Account bool `flag:"account" desc:"Refresh the account cache"`
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the cache",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[refreshFlags](cmd)
		if err != nil {
			return err
		}
		return refreshCaches(flags.Account)
	},
}

func refreshCaches(account bool) error {
	if account {
		_, err := api.MyDetails()
		if err != nil {
			return err
		}
	}
	cache.CleanBank()
	cache.CleanBankItems()
	cache.CleanCharacters()
	_, err := api.MyBank()
	if err != nil {
		return err
	}
	_, err = api.MyBankItems()
	if err != nil {
		return err
	}
	_, err = api.AccountsCharacters("")
	return err
}

func init() {
	rootCmd.AddCommand(refreshCmd)
	err := utils.RegisterFlags[refreshFlags](refreshCmd)
	if err != nil {
		panic(err)
	}
}
