package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var npcsDetailsCmd = &cobra.Command{
	Use:   "npcsDetails [code]",
	Short: "Get All Npcs",
	Long: `Get All Npcs

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
			path = "/npcs/details"
		case 1:
			path = fmt.Sprintf("/npcs/details/%s", args[0])
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
	apiCmd.AddCommand(npcsDetailsCmd)
	npcsDetailsCmd.Flags().String("currency", "",
		"Currency code to filter NPCs that trade with this currency.")
	npcsDetailsCmd.Flags().String("item", "",
		"Item code to filter NPCs that trade this item.")
	npcsDetailsCmd.Flags().String("name", "",
		"NPC name.")
	npcsDetailsCmd.Flags().Int("page", 0,
		"Page number")
	npcsDetailsCmd.Flags().Int("size", 0,
		"Page size")
	npcsDetailsCmd.Flags().String("type", "",
		"Type of NPCs.")
}
