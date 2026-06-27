package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var charactersStatsCmd = &cobra.Command{
	Use:   "charactersStats <name>",
	Short: "Get Character Stats",
	Long: `Get Character Stats

Arguments:
  name   The name of the character.`,
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
			path = fmt.Sprintf("/characters/%s/stats", args[0])
		}

		params := make(map[string]string)

		resp, err := api.Get(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(charactersStatsCmd)
}
