package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var accountsAchievementsCmd = &cobra.Command{
	Use:   "accountsAchievements [account]",
	Short: "Get Account Achievements",
	Long: `Get Account Achievements

Arguments:
  [account]  The name of the account.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			return fmt.Errorf("missing required argument: account")
		}
		path = fmt.Sprintf("/accounts/%s/achievements", args[0])

		params := make(map[string]string)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			params[f.Name] = f.Value.String()
		})

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return output(resp)
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
