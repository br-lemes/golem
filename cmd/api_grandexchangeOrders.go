package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var grandexchangeOrdersCmd = &cobra.Command{
	Use:   "grandexchangeOrders [id]",
	Short: "Get Ge Orders",
	Long: `Get Ge Orders

Arguments:
  [id]  The id of the order.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			path = "/grandexchange/orders"
		} else if len(args) == 1 {
			path = fmt.Sprintf("/grandexchange/orders/%s", args[0])
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
	apiCmd.AddCommand(grandexchangeOrdersCmd)
	grandexchangeOrdersCmd.Flags().String("account", "",
		"The account that sells or buys items.")
	grandexchangeOrdersCmd.Flags().String("code", "",
		"The code of the item.")
	grandexchangeOrdersCmd.Flags().Int("page", 0,
		"Page number")
	grandexchangeOrdersCmd.Flags().Int("size", 0,
		"Page size")
	grandexchangeOrdersCmd.Flags().String("type", "",
		"Filter by order type (sell or buy).")
}
