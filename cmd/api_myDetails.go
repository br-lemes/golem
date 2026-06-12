package cmd

import (
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var myDetailsCmd = &cobra.Command{
	Use:   "myDetails",
	Short: "Get Account Details",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/my/details"

		params := make(map[string]string)

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(myDetailsCmd)
}
