package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/br-lemes/golem/pkg/schemas"
)

var (
	//go:embed resources.json
	resources []byte

	resourcesList  []ResourceSchema
	resourcesCache map[string]ResourceSchema
)

func GetResources() []ResourceSchema {
	initResourcesCache()
	return resourcesList
}

func GetResource(code string) (ResourceSchema, bool) {
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

func initResourcesCache() {
	if resourcesCache != nil {
		return
	}
	resourcesCache = make(map[string]ResourceSchema)
	err := json.Unmarshal(resources, &resourcesList)
	if err != nil {
		panic("failed to unmarshal resources: " + err.Error())
	}
	for _, resource := range resourcesList {
		resourcesCache[resource.Code] = resource
	}
}
