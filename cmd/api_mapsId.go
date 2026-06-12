package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var mapsIdCmd = &cobra.Command{
	Use:   "mapsId [map_id]",
	Short: "Get Map By Id",
	Long: `Get Map By Id

Arguments:
  [map_id]  The unique ID of the map.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			return fmt.Errorf("missing required argument: map_id")
		}
		path = fmt.Sprintf("/maps/id/%s", args[0])

		params := make(map[string]string)

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(mapsIdCmd)
}
