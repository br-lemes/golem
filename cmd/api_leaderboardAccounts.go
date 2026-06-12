package cmd

import (
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var leaderboardAccountsCmd = &cobra.Command{
	Use:   "leaderboardAccounts",
	Short: "Get Accounts Leaderboard",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/leaderboard/accounts"

		params := make(map[string]string)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			params[f.Name] = f.Value.String()
		})

		resp, err := apiGet(path, params)
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
