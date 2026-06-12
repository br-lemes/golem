package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var myActionFightCmd = &cobra.Command{
	Use:   "myActionFight <name> [participants...]",
	Short: "Action Fight",
	Long: `Action Fight

Arguments:
  name           Name of your character.
  participants   Names of other characters in your account to fight alongside.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 1, 2, 3:
			return nil
		default:
			return fmt.Errorf("invalid number of arguments: %d", argCount)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		argCount := len(args)

		path = fmt.Sprintf("/my/%s/action/fight", args[0])

		participants := []string{}
		if argCount > 1 {
			participants = args[1:]
		}

		params := map[string][]string{
			"participants": participants,
		}

		resp, err := api.Post(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(myActionFightCmd)
}
