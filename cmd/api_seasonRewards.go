package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var seasonRewardsCmd = &cobra.Command{
	Use:   "seasonRewards [code]",
	Short: "Get All Season Rewards",
	Long: `Get All Season Rewards

Arguments:
  code   The code of the season reward.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 0:
			return nil
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
		case 0:
			path = "/season_rewards"
		case 1:
			path = fmt.Sprintf("/season_rewards/%s", args[0])
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
	apiCmd.AddCommand(seasonRewardsCmd)
	seasonRewardsCmd.Flags().Int("page", 0, "Page number")
	seasonRewardsCmd.Flags().Int("size", 0, "Page size")
	seasonRewardsCmd.Flags().String("type", "", "Filter by reward type.")
}
