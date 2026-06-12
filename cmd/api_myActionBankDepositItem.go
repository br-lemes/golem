package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var myActionBankDepositItemCmd = &cobra.Command{
	Use:   "myActionBankDepositItem <name>",
	Short: "Action Deposit Bank Item",
	Long: `Action Deposit Bank Item

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
			path = fmt.Sprintf("/my/%s/action/bank/deposit/item", args[0])
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
			params[f.Name] = f.Value.String()
		})

		resp, err := apiPost(path, []map[string]string{params})
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionBankDepositItemCmd)
	myActionBankDepositItemCmd.Flags().String("code", "",
		"Item code.")
	myActionBankDepositItemCmd.Flags().Int("quantity", 0,
		"Item quantity.")
}
