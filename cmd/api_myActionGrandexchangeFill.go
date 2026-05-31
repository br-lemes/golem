package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var myActionGrandexchangeFillCmd = &cobra.Command{
	Use:   "myActionGrandexchangeFill <name>",
	Short: "Action Ge Fill",
	Long: `Action Ge Fill

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
			path = fmt.Sprintf("/my/%s/action/grandexchange/fill", args[0])
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
	apiCmd.AddCommand(myActionGrandexchangeFillCmd)
	myActionGrandexchangeFillCmd.Flags().String("id", "",
		"Buy order id.")
	myActionGrandexchangeFillCmd.Flags().Int("quantity", 0,
		"Item quantity to sell.")
}
