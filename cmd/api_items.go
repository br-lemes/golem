package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var itemsCmd = &cobra.Command{
	Use:   "items [code]",
	Short: "Get All Items",
	Long: `Get All Items

Arguments:
  code   The code of the item.`,
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
			path = "/items"
		case 1:
			path = fmt.Sprintf("/items/%s", args[0])
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
	apiCmd.AddCommand(itemsCmd)
	itemsCmd.Flags().String("craft_material", "", "Item code of items used as material for crafting.")
	itemsCmd.Flags().String("craft_skill", "", "Skill to craft items.")
	itemsCmd.Flags().Int("max_level", 0, "Maximum level.")
	itemsCmd.Flags().Int("min_level", 0, "Minimum level.")
	itemsCmd.Flags().String("name", "", "Name of the item.")
	itemsCmd.Flags().Int("page", 0, "Page number")
	itemsCmd.Flags().Int("size", 0, "Page size")
	itemsCmd.Flags().String("type", "", "Type of items.")
}
