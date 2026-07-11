package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var accountsAchievementsCmd = &cobra.Command{
	Use:   "accountsAchievements <account>",
	Short: "Get Account Achievements",
	Long: `Get Account Achievements

Arguments:
  account   The name of the account.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 1:
			return nil
		default:
			return fmt.Errorf("invalid number of arguments: %d", argCount)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		argCount := len(args)

		switch argCount {
		case 1:
			path = fmt.Sprintf("/accounts/%s/achievements", args[0])
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
	apiCmd.AddCommand(accountsAchievementsCmd)
	accountsAchievementsCmd.Flags().Bool("completed", false,
		"Filter by completed achievements.")
	accountsAchievementsCmd.Flags().Int("page", 0,
		"Page number")
	accountsAchievementsCmd.Flags().Int("size", 0,
		"Page size")
	accountsAchievementsCmd.Flags().String("type", "",
		"Type of achievements.")
}
