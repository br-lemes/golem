package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var myActionBankBuyExpansionCmd = &cobra.Command{
	Use:   "myActionBankBuyExpansion <name>",
	Short: "Action Buy Bank Expansion",
	Long: `Action Buy Bank Expansion

Arguments:
  name   Name of your character.`,
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
			path = fmt.Sprintf("/my/%s/action/bank/buy_expansion", args[0])
		}

		params := make(map[string]string)

		resp, err := apiPost(path, params)
		if err != nil {
			return err
		}
		return output(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionBankBuyExpansionCmd)
}
