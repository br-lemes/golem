package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var myActionTaskCancelCmd = &cobra.Command{
	Use:   "myActionTaskCancel <name>",
	Short: "Action Task Cancel",
	Long: `Action Task Cancel

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
			path = fmt.Sprintf("/my/%s/action/task/cancel", args[0])
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
	apiCmd.AddCommand(myActionTaskCancelCmd)
}
