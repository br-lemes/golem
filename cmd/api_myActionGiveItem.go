package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var myActionGiveItemCmd = &cobra.Command{
	Use:   "myActionGiveItem <name>",
	Short: "Action Give Items",
	Long: `Action Give Items

Arguments:
  name   Name of your character.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
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
		case 1:
			path = fmt.Sprintf("/my/%s/action/give/item", args[0])
		}

		character, _ := cmd.Flags().GetString("character")
		code, _ := cmd.Flags().GetString("code")
		quantity, _ := cmd.Flags().GetInt("quantity")

		params := map[string]interface{}{
			"character": character,
			"items": []map[string]interface{}{
				{
					"code":     code,
					"quantity": quantity,
				},
			},
		}

		resp, err := apiPost(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionGiveItemCmd)
	myActionGiveItemCmd.Flags().String("character", "",
		"Character name. The name of the character who will receive the items.")
	myActionGiveItemCmd.Flags().String("code", "",
		"Item code.")
	myActionGiveItemCmd.Flags().Int("quantity", 0,
		"Item quantity.")
}
