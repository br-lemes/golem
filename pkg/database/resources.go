package database

import (
	_ "embed"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed resources.json
var resources []byte

var Resources = newStore(jsonLoader[schemas.ResourceSchema](resources), func(resource *schemas.ResourceSchema) string {
	return resource.Code
})
