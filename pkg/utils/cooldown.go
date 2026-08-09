package utils

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/getkin/kin-openapi/openapi3"
)

func hasCooldown(pathItem *openapi3.PathItem) bool {
	for _, o := range pathItem.Operations() {
		if o.Responses != nil && o.Responses.Status(499) != nil {
			return true
		}
	}
	return false
}

func GetCooldown(targetPath string) (bool, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(database.OpenAPI())
	if err != nil {
		return false, err
	}

	pathItem := doc.Paths.Find(targetPath)
	if pathItem == nil {
		return false, fmt.Errorf("path %s not found", targetPath)
	}

	return hasCooldown(pathItem), nil
}

func GetCooldowns() (map[string]bool, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(database.OpenAPI())
	if err != nil {
		return nil, err
	}

	pathsMap := doc.Paths.Map()
	result := make(map[string]bool)

	for path, pathItem := range pathsMap {
		if hasCooldown(pathItem) {
			result[path] = true
		}
	}

	return result, nil
}
