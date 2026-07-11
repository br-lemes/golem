package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var grandexchangeHistoryCmd = &cobra.Command{
	Use:   "grandexchangeHistory <code>",
	Short: "Get Ge History",
	Long: `Get Ge History

Arguments:
  code   The code of the item.`,
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
			path = fmt.Sprintf("/grandexchange/history/%s", args[0])
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
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
	apiCmd.AddCommand(grandexchangeHistoryCmd)
	grandexchangeHistoryCmd.Flags().String("account", "",
		"Account involved in the transaction (matches either seller or buyer).")
	grandexchangeHistoryCmd.Flags().Int("page", 0,
		"Page number")
	grandexchangeHistoryCmd.Flags().Int("size", 0,
		"Page size")
}
