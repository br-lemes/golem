package cmd

import (
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var myPendingItemsCmd = &cobra.Command{
	Use:   "myPendingItems",
	Short: "Get Pending Items",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/my/pending-items"

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
	apiCmd.AddCommand(myPendingItemsCmd)
	myPendingItemsCmd.Flags().Int("page", 0,
		"Page number")
	myPendingItemsCmd.Flags().Int("size", 0,
		"Page size")
}
