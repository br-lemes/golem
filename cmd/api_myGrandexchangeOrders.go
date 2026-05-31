package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var myGrandexchangeOrdersCmd = &cobra.Command{
	Use:   "myGrandexchangeOrders",
	Short: "Get Ge Orders",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/my/grandexchange/orders"

		params := make(map[string]string)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			params[f.Name] = f.Value.String()
		})

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return output(resp)
	},
}

func init() {
	apiCmd.AddCommand(myGrandexchangeOrdersCmd)
	myGrandexchangeOrdersCmd.Flags().String("code", "",
		"The code of the item.")
	myGrandexchangeOrdersCmd.Flags().Int("page", 0,
		"Page number")
	myGrandexchangeOrdersCmd.Flags().Int("size", 0,
		"Page size")
	myGrandexchangeOrdersCmd.Flags().String("type", "",
		"Filter by order type (sell or buy).")
}
