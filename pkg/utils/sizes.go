package utils

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/getkin/kin-openapi/openapi3"
)

func sizeMax(operation *openapi3.Operation) (int, bool) {
	for _, paramRef := range operation.Parameters {
		param := paramRef.Value
		if param == nil || param.Name != "size" {
			continue
		}
		if param.Schema != nil && param.Schema.Value != nil && param.Schema.Value.Max != nil {
			return int(*param.Schema.Value.Max), true
		}
	}
	return 0, false
}

func GetSize(targetPath string) (int, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(database.OpenAPI())
	if err != nil {
		return 0, err
	}

	pathItem := doc.Paths.Find(targetPath)
	if pathItem == nil {
		return 0, fmt.Errorf("path %s not found", targetPath)
	}

	found := false
	result := 0
	for method, operation := range pathItem.Operations() {
		max, ok := sizeMax(operation)
		if !ok {
			continue
		}
		if found {
			return 0, fmt.Errorf("multiple methods with size parameter for path %s (at least %s)", targetPath, method)
		}
		found = true
		result = max
	}

	if !found {
		return 0, fmt.Errorf("size parameter not found for path %s", targetPath)
	}

	return result, nil
}

func GetSizes() (map[string]int, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(database.OpenAPI())
	if err != nil {
		return nil, err
	}

	pathsMap := doc.Paths.Map()
	result := make(map[string]int)

	for path, pathItem := range pathsMap {
		found := false
		for method, operation := range pathItem.Operations() {
			max, ok := sizeMax(operation)
			if !ok {
				continue
			}
			if found {
				return nil, fmt.Errorf("multiple methods with size parameter for path %s (at least %s)", path, method)
			}
			found = true
			result[path] = max
		}
	}

	return result, nil
}
