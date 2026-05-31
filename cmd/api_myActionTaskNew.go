package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var myActionTaskNewCmd = &cobra.Command{
	Use:   "myActionTaskNew <name>",
	Short: "Action Accept New Task",
	Long: `Action Accept New Task

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
			path = fmt.Sprintf("/my/%s/action/task/new", args[0])
		}

		params := make(map[string]string)

		resp, err := apiPost(path, params)
		if err != nil {
			return err
		}
		return output(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionTaskNewCmd)
}
