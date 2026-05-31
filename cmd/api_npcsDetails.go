package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var npcsDetailsCmd = &cobra.Command{
	Use:   "npcsDetails [code]",
	Short: "Get All Npcs",
	Long: `Get All Npcs

Arguments:
  [code]  The code of the NPC.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			path = "/npcs/details"
		} else if len(args) == 1 {
			path = fmt.Sprintf("/npcs/details/%s", args[0])
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
