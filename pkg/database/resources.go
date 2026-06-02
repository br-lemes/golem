package database

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/br-lemes/golem/pkg/schemas"
)

//go:embed resources.json
var resources []byte

var resourcesCache map[string]ResourceSchema

func GetResource(code string) (ResourceSchema, bool) {
	sync.OnceFunc(initResourcesCache)()
	resource, exists := resourcesCache[code]
	return resource, exists
}

func initResourcesCache() {
	if resourcesCache != nil {
		return
	}
	resourcesCache = make(map[string]ResourceSchema)
	var result []ResourceSchema
	err := json.Unmarshal(resources, &result)
	if err != nil {
		return
	}
	for _, resource := range result {
		resourcesCache[resource.Code] = resource
	}
}
