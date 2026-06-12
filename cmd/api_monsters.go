package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var monstersCmd = &cobra.Command{
	Use:   "monsters [code]",
	Short: "Get All Monsters",
	Long: `Get All Monsters

Arguments:
  code   The code of the monster.`,
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
			path = "/monsters"
		case 1:
			path = fmt.Sprintf("/monsters/%s", args[0])
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
	apiCmd.AddCommand(monstersCmd)
	monstersCmd.Flags().String("drop", "",
		"Item code of the drop.")
	monstersCmd.Flags().Int("max_level", 0,
		"Maximum level.")
	monstersCmd.Flags().Int("min_level", 0,
		"Minimum level.")
	monstersCmd.Flags().String("name", "",
		"Name of the monster.")
	monstersCmd.Flags().Int("page", 0,
		"Page number")
	monstersCmd.Flags().Int("size", 0,
		"Page size")
}
