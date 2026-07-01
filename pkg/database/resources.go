package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/tidwall/gjson"
)

var (
	//go:embed resources.json
	resources []byte

	resourcesList     []schemas.ResourceSchema
	resourcesCache    map[string]schemas.ResourceSchema
	resourceCodesList []string
)

func GetResources() []schemas.ResourceSchema {
	initResourcesCache()
	return resourcesList
}

func GetResource(code string) (schemas.ResourceSchema, bool) {
	initResourcesCache()
	resource, exists := resourcesCache[code]
	return resource, exists
}

func GetResourceCodes() []string {
	initResourceCodes()
	return resourceCodesList
}

var initResourcesCache = sync.OnceFunc(func() {
	resourcesCache = make(map[string]schemas.ResourceSchema)
	err := json.Unmarshal(resources, &resourcesList)
	if err != nil {
		panic("failed to unmarshal resources: " + err.Error())
	}
	for _, resource := range resourcesList {
		resourcesCache[resource.Code] = resource
	}
})

var initResourceCodes = sync.OnceFunc(func() {
	gjson.GetBytes(resources, "#.code").ForEach(func(_, value gjson.Result) bool {
		resourceCodesList = append(resourceCodesList, value.String())
		return true
	})
})
