package cmd

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var commandsCmd = &cobra.Command{
	Use:   "commands [command]",
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
			return console.Auto(buildCompactMap(command))
		}
		return console.Auto(commands)
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

func buildCompactMap(routes []map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for _, route := range routes {
		for method, path := range route {
			returnType := fetchReturnTypeFromSpec(method, path)

			_, exists := result[method]
			if !exists {
				result[method] = make(map[string]string)
			}

			result[method][path] = returnType
		}
	}
	return result
}

func fetchReturnTypeFromSpec(method string, targetPath string) string {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(database.OpenAPI())
	if err != nil {
		return ""
	}

	pathItem := doc.Paths.Find(targetPath)
	if pathItem == nil {
		return ""
	}

	upperMethod := strings.ToUpper(method)
	operation := pathItem.GetOperation(upperMethod)
	if operation == nil {
		return ""
	}

	responseRef := operation.Responses.Status(200)
	if responseRef == nil {
		return "void"
	}

	response := responseRef.Value
	if response == nil {
		return "void"
	}

	jsonContent := response.Content.Get("application/json")
	if jsonContent == nil {
		return "void"
	}

	if jsonContent.Schema == nil {
		return "void"
	}

	schemaRef := jsonContent.Schema
	if schemaRef.Ref != "" {
		return path.Base(schemaRef.Ref)
	}

	return "void"
}
