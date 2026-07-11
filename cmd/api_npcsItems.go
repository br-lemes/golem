package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var npcsItemsCmd = &cobra.Command{
	Use:   "npcsItems [code]",
	Short: "Get All Npcs Items",
	Long: `Get All Npcs Items

Arguments:
  code   The code of the NPC.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 0:
			return nil
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
		case 0:
			path = "/npcs/items"
		case 1:
			path = fmt.Sprintf("/npcs/items/%s", args[0])
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
	apiCmd.AddCommand(npcsItemsCmd)
	npcsItemsCmd.Flags().String("code", "",
		"Item code.")
	npcsItemsCmd.Flags().String("currency", "",
		"Currency code.")
	npcsItemsCmd.Flags().String("npc", "",
		"NPC code.")
	npcsItemsCmd.Flags().Int("page", 0,
		"Page number")
	npcsItemsCmd.Flags().Int("size", 0,
		"Page size")
}
