package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var commandsCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "commands [command]",
	Short: "List available commands mapped from the game's OpenAPI spec",
	Long: `List available commands mapped from the game's OpenAPI spec

Arguments:
  command   Name of a specific command to inspect.`,
	ValidArgsFunction: completion.Custom(1, utils.GetCommandsCompletion).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		commands, err := utils.GetCommands()
		if err != nil {
			return err
		}

		if len(args) > 0 {
			target := args[0]
			_, exists := commands[target]
			if !exists {
				return fmt.Errorf("command not found: %s", target)
			}
		}

		missing, _ := cmd.Flags().GetBool("missing")
		if missing {
			for _, command := range apiCmd.Commands() {
				delete(commands, command.Name())
			}
		}

		if len(args) > 0 {
			command := commands[args[0]]
			if command == nil {
				return nil
			}
			return console.Auto(utils.BuildCompactMap(command))
		}
		return console.Auto(commands)
	},
}

func init() {
	rootCmd.AddCommand(commandsCmd)
	commandsCmd.Flags().BoolP("missing", "m", false, "Show only commands with missing implementation")
}
