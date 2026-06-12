package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var achievementsCmd = &cobra.Command{
	Use:   "achievements [code]",
	Short: "Get All Achievements",
	Long: `Get All Achievements

Arguments:
  [code]  The code of the achievement.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			path = "/achievements"
		} else if len(args) == 1 {
			path = fmt.Sprintf("/achievements/%s", args[0])
		}

		params := make(map[string]string)
		cmd.Flags().Visit(func(f *pflag.Flag) {
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
	apiCmd.AddCommand(achievementsCmd)
	achievementsCmd.Flags().Int("page", 0,
		"Page number")
	achievementsCmd.Flags().Int("size", 0,
		"Page size")
	achievementsCmd.Flags().String("type", "",
		"Type of achievements.")
}
