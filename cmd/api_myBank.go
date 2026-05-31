package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var myBankCmd = &cobra.Command{
	Use:   "myBank",
	Short: "Get Bank Details",
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
			path = "/my/bank"
		}

		params := make(map[string]string)

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return output(resp)
	},
}

func init() {
	apiCmd.AddCommand(myBankCmd)
}
