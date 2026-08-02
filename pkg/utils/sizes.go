package utils

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/getkin/kin-openapi/openapi3"
)

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

	operation := pathItem.Get
	if operation == nil {
		return 0, fmt.Errorf("GET method not found for path %s", targetPath)
	}

	for _, paramRef := range operation.Parameters {
		param := paramRef.Value
		if param == nil || param.Name != "size" {
			continue
		}
		if param.Schema != nil && param.Schema.Value != nil &&
			param.Schema.Value.Max != nil {
			return int(*param.Schema.Value.Max), nil
		}
		return 0, fmt.Errorf("size parameter has no maximum for path %s",
			targetPath)
	}

	return 0, fmt.Errorf("size parameter not found for path %s", targetPath)
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
		operation := pathItem.Get
		if operation == nil {
			continue
		}

		for _, paramRef := range operation.Parameters {
			param := paramRef.Value
			if param != nil && param.Name == "size" &&
				param.Schema != nil && param.Schema.Value != nil &&
				param.Schema.Value.Max != nil {
				result[path] = int(*param.Schema.Value.Max)
				break
			}
		}
	}

	return result, nil
}
