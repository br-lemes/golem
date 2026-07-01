package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/br-lemes/golem/pkg/schemas"
)

var (
	//go:embed resources.json
	resources []byte

	resourcesList  []schemas.ResourceSchema
	resourcesCache map[string]schemas.ResourceSchema
)

func GetResources() []schemas.ResourceSchema {
	initResourcesCache()
	return resourcesList
}

func GetResource(code string) (schemas.ResourceSchema, bool) {
	sync.OnceFunc(initResourcesCache)()
	resource, exists := resourcesCache[code]
	return resource, exists
}

func GetResourceCodes() []string {
	initResourcesCache()
	codes := make([]string, 0, len(resourcesCache))
	for code := range resourcesCache {
		codes = append(codes, code)
	}
	return codes
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
