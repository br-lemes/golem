package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var myRatesCmd = &cobra.Command{
	Use:   "myRates",
	Short: "Get Rate Limits",
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
			path = "/my/rates"
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
	apiCmd.AddCommand(myRatesCmd)
}
