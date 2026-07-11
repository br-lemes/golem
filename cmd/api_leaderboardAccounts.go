package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var leaderboardAccountsCmd = &cobra.Command{
	Use:   "leaderboardAccounts",
	Short: "Get Accounts Leaderboard",
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 0:
			return nil
		default:
			return fmt.Errorf("invalid number of arguments: %d", argCount)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		argCount := len(args)

		switch argCount {
		case 0:
			path = "/leaderboard/accounts"
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
			params[f.Name] = f.Value.String()
		})

		resp, err := api.Get(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(leaderboardAccountsCmd)
	leaderboardAccountsCmd.Flags().String("name", "",
		"Account name.")
	leaderboardAccountsCmd.Flags().Int("page", 0,
		"Page number")
	leaderboardAccountsCmd.Flags().Int("size", 0,
		"Page size")
	leaderboardAccountsCmd.Flags().String("sort", "",
		"Sort of account leaderboards.")
}
