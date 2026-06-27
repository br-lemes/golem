package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var sandboxGiveXpCmd = &cobra.Command{
	Use:   "sandboxGiveXp",
	Short: "Give Xp",
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
			path = "/sandbox/give_xp"
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
			params[f.Name] = f.Value.String()
		})

		resp, err := api.Post(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(sandboxGiveXpCmd)
	sandboxGiveXpCmd.Flags().Int("amount", 0,
		"Amount of XP to give to the character.")
	sandboxGiveXpCmd.Flags().String("character", "",
		"Name of the character to receive the XP.")
	sandboxGiveXpCmd.Flags().String("type", "",
		"Type of XP to give (e.g., 'combat', 'woodcutting', 'mining', etc.).")
}
