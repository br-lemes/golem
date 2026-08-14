package utils

import (
	"fmt"
	"sort"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/getkin/kin-openapi/openapi3"
)

type FieldInfo struct {
	Name        string `json:"name"`
	In          string `json:"in,omitempty"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type MethodDetails struct {
	Summary     string      `json:"summary"`
	Parameters  []FieldInfo `json:"parameters,omitempty"`
	RequestBody []FieldInfo `json:"requestBody,omitempty"`
	IsBodyArray bool        `json:"isBodyArray,omitempty"`
}

type RouteData map[string]MethodDetails

func GetRoute(targetPath string) (RouteData, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(database.OpenAPI())
	if err != nil {
		return nil, err
	}

	pathItem := doc.Paths.Find(targetPath)
	if pathItem == nil {
		return nil, fmt.Errorf("path %s not found", targetPath)
	}

	result := make(RouteData)
	for method, operation := range pathItem.Operations() {
		var queryParams []FieldInfo

		for _, paramRef := range operation.Parameters {
			param := paramRef.Value
			if param != nil {
				paramType := ""
				if param.Schema != nil && param.Schema.Value != nil && param.Schema.Value.Type != nil {
					types := *param.Schema.Value.Type
					if len(types) > 0 {
						paramType = types[0]
					}
				}

				queryParams = append(queryParams, FieldInfo{
					Name:        param.Name,
					In:          param.In,
					Description: param.Description,
					Type:        paramType,
				})
			}
		}

		sort.Slice(queryParams, func(i, j int) bool {
			return queryParams[i].Name < queryParams[j].Name
		})

		var bodyFields []FieldInfo
		isBodyArray := false

		if operation.RequestBody != nil && operation.RequestBody.Value != nil {
			reqBody := operation.RequestBody.Value
			jsonContent := reqBody.Content.Get("application/json")
			if jsonContent != nil && jsonContent.Schema != nil && jsonContent.Schema.Value != nil {

				schema := jsonContent.Schema.Value
				schemaType := ""
				if schema.Type != nil && len(*schema.Type) > 0 {
					schemaType = (*schema.Type)[0]
				}

				if schemaType == "array" {
					isBodyArray = true
					if schema.Items != nil && schema.Items.Value != nil {
						itemSchema := schema.Items.Value
						for propName, propRef := range itemSchema.Properties {
							prop := propRef.Value
							if prop != nil {
								propType := ""
								if prop.Type != nil && len(*prop.Type) > 0 {
									propType = (*prop.Type)[0]
								}

								bodyFields = append(bodyFields, FieldInfo{
									Name:        propName,
									Description: prop.Description,
									Type:        propType,
								})
							}
						}
					}
				} else {
					for propName, propRef := range schema.Properties {
						prop := propRef.Value
						if prop != nil {
							propType := ""
							if prop.Type != nil && len(*prop.Type) > 0 {
								propType = (*prop.Type)[0]
							}

							bodyFields = append(bodyFields, FieldInfo{
								Name:        propName,
								Description: prop.Description,
								Type:        propType,
							})
						}
					}
				}
			}
		}

		sort.Slice(bodyFields, func(i, j int) bool {
			return bodyFields[i].Name < bodyFields[j].Name
		})

		methodData := MethodDetails{
			Summary:     operation.Summary,
			IsBodyArray: isBodyArray,
		}

		if len(queryParams) > 0 {
			methodData.Parameters = queryParams
		}

		if len(bodyFields) > 0 {
			methodData.RequestBody = bodyFields
		}

		result[method] = methodData
	}

	return result, nil
}

func GetRoutes() ([]map[string]string, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(database.OpenAPI())
	if err != nil {
		return nil, err
	}

	pathsMap := doc.Paths.Map()
	paths := make([]string, 0, len(pathsMap))
	for path := range pathsMap {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var result []map[string]string
	for _, path := range paths {
		pathItem := pathsMap[path]
		for method := range pathItem.Operations() {
			result = append(result, map[string]string{method: path})
		}
	}

	return result, nil
}

func GetRoutesCompletion() []string {
	paths, _ := GetRoutes()
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		for method := range path {
			result = append(result, fmt.Sprintf("%s\t%s", path[method], method))
		}
	}
	return result
}
