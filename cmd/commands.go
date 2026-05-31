package cmd

import (
	"fmt"
	"regexp"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var commandsCmd = &cobra.Command{
	Use:   "commands <command>",
	Short: "List available commands mapped from the game's OpenAPI spec",
	Long: `List available commands mapped from the game's OpenAPI spec

Arguments:
  command   Name of a specific command to inspect.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		commands, err := getCommands()
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
			return output(command)
		}
		return output(commands)
	},
}

func init() {
	rootCmd.AddCommand(commandsCmd)
	commandsCmd.Flags().BoolP("missing", "m", false,
		"Show only commands with missing implementation")
}

func getCommands() (map[string][]map[string]string, error) {
	routes, err := getRoutes()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\{[^}]*\}|/`)
	result := make(map[string][]map[string]string)
	for _, route := range routes {
		for method, path := range route {
			key := strcase.ToLowerCamel(re.ReplaceAllString(path, " "))
			result[key] = append(result[key], map[string]string{method: path})
		}
	}

	return result, nil
}
