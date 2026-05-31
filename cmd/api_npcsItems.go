package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var npcsItemsCmd = &cobra.Command{
	Use:   "npcsItems [code]",
	Short: "Get All Npcs Items",
	Long: `Get All Npcs Items

Arguments:
  [code]  The code of the NPC.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			path = "/npcs/items"
		} else if len(args) == 1 {
			path = fmt.Sprintf("/npcs/items/%s", args[0])
		}

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
