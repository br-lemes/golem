package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var myGrandexchangeHistoryCmd = &cobra.Command{
	Use:   "myGrandexchangeHistory",
	Short: "Get Ge History",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/my/grandexchange/history"

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
	apiCmd.AddCommand(myGrandexchangeHistoryCmd)
	myGrandexchangeHistoryCmd.Flags().String("code", "",
		"Item to search in your history.")
	myGrandexchangeHistoryCmd.Flags().String("id", "",
		"Order ID to search in your history.")
	myGrandexchangeHistoryCmd.Flags().Int("page", 0,
		"Page number")
	myGrandexchangeHistoryCmd.Flags().Int("size", 0,
		"Page size")
}
