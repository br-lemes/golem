package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var grandexchangeHistoryCmd = &cobra.Command{
	Use:   "grandexchangeHistory [code]",
	Short: "Get Ge History",
	Long: `Get Ge History

Arguments:
  [code]  The code of the item.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			return fmt.Errorf("missing required argument: code")
		}
		path = fmt.Sprintf("/grandexchange/history/%s", args[0])

		params := make(map[string]string)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			params[f.Name] = f.Value.String()
		})

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(grandexchangeHistoryCmd)
	grandexchangeHistoryCmd.Flags().String("account", "",
		"Account involved in the transaction (matches either seller or buyer).")
	grandexchangeHistoryCmd.Flags().Int("page", 0,
		"Page number")
	grandexchangeHistoryCmd.Flags().Int("size", 0,
		"Page size")
}
