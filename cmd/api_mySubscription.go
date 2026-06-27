package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var mySubscriptionCmd = &cobra.Command{
	Use:   "mySubscription",
	Short: "Get My Subscription",
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
			path = "/my/subscription"
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
	apiCmd.AddCommand(mySubscriptionCmd)
}
