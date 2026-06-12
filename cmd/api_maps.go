package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var mapsCmd = &cobra.Command{
	Use:   "maps [layer | layer x y]",
	Short: "Get All Maps",
	Long: `Get All Maps

Arguments:
  layer   The layer of the map (interior, overworld, underground).
  x       The position x of the map.
  y       The position y of the map.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 0:
			return nil
		case 1:
			return nil
		case 3:
			return nil
		default:
			return fmt.Errorf("invalid number of arguments: %d", argCount)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		argCount := len(args)

		switch argCount {
		case 0:
			path = "/maps"
		case 1:
			path = fmt.Sprintf("/maps/%s", args[0])
		case 3:
			path = fmt.Sprintf("/maps/%s/%s/%s", args[0], args[1], args[2])
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
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
	apiCmd.AddCommand(mapsCmd)
	mapsCmd.Flags().String("content_code", "",
		"Content code on the map.")
	mapsCmd.Flags().String("content_type", "",
		"Type of maps.")
	mapsCmd.Flags().Bool("hide_blocked_maps", false,
		"When true, excludes maps with access_type 'blocked' from the results.")
	mapsCmd.Flags().String("layer", "",
		"Filter maps by layer.")
	mapsCmd.Flags().Int("page", 0,
		"Page number")
	mapsCmd.Flags().Int("size", 0,
		"Page size")
}
