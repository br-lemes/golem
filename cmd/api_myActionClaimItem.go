package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var myActionClaimItemCmd = &cobra.Command{
	Use:   "myActionClaimItem <name id>",
	Short: "Action Claim Pending Item",
	Long: `Action Claim Pending Item

Arguments:
  name   Name of your character.
  id     The ID of the pending item to claim.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 2:
			return nil
		default:
			return fmt.Errorf("invalid number of arguments: %d", argCount)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		argCount := len(args)

		switch argCount {
		case 2:
			path = fmt.Sprintf("/my/%s/action/claim_item/%s", args[0], args[1])
		}

		params := make(map[string]string)

		resp, err := api.Post(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionClaimItemCmd)
}
