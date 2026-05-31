package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var myActionUseCmd = &cobra.Command{
	Use:   "myActionUse <name>",
	Short: "Action Use Item",
	Long: `Action Use Item

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
			path = fmt.Sprintf("/my/%s/action/use", args[0])
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
			params[f.Name] = f.Value.String()
		})

		resp, err := apiPost(path, params)
		if err != nil {
			return err
		}
		return output(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionUseCmd)
	myActionUseCmd.Flags().String("code", "",
		"Item code.")
	myActionUseCmd.Flags().Int("quantity", 0,
		"Item quantity.")
}
