package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var charactersCmd = &cobra.Command{
	Use:   "characters [name]",
	Short: "Get Character",
	Long: `Get Character

Arguments:
  [name]  The name of the character.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			return fmt.Errorf("missing required argument: name")
		}
		path = fmt.Sprintf("/characters/%s", args[0])

		params := make(map[string]string)

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return output(resp)
	},
}

func init() {
	apiCmd.AddCommand(charactersCmd)
}
