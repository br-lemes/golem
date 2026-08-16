// +gocover:ignore:file only returns embedded data
package database

import _ "embed"

//go:embed openapi.json
var openapi []byte

func OpenAPI() []byte {
	return openapi
}
