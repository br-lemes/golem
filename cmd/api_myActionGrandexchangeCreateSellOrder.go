package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var myActionGrandexchangeCreateSellOrderCmd = &cobra.Command{
	Use:   "myActionGrandexchangeCreateSellOrder <name>",
	Short: "Action Ge Create Sell Order",
	Long: `Action Ge Create Sell Order

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
			path = fmt.Sprintf("/my/%s/action/grandexchange/create_sell_order", args[0])
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
			params[f.Name] = f.Value.String()
		})

		resp, err := api.Post(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionGrandexchangeCreateSellOrderCmd)
	myActionGrandexchangeCreateSellOrderCmd.Flags().String("code", "",
		"Item code.")
	myActionGrandexchangeCreateSellOrderCmd.Flags().Int("price", 0,
		"Item price per unit.")
	myActionGrandexchangeCreateSellOrderCmd.Flags().Int("quantity", 0,
		"Item quantity.")
}
