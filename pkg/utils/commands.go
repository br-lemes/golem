package utils

import (
	"path"
	"regexp"
	"strings"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

func GetCommands() (map[string][]map[string]string, error) {
	routes, err := GetRoutes()
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

func GetCommandsCompletion() []string {
	commands, _ := GetCommands()
	result := make([]string, 0, len(commands))
	for command := range commands {
		result = append(result, command)
	}
	return result
}

func BuildCompactMap(routes []map[string]string) map[string]map[string]string {
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
