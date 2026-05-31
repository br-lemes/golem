package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var charactersCreateCmd = &cobra.Command{
	Use:   "charactersCreate",
	Short: "Create Character",
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 0:
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
			path = "/characters/create"
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
	apiCmd.AddCommand(charactersCreateCmd)
	charactersCreateCmd.Flags().String("name", "",
		"Your desired character name. It's unique and all players can see it.")
	charactersCreateCmd.Flags().String("skin", "",
		"Your desired skin. Skins unlocked by default: 'men1', 'men2', 'men3', 'women1', 'women2', 'women3'.")
}
