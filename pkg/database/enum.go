package database

import (
	"sync"

	"github.com/tidwall/gjson"
)

type enumData struct {
	values map[string][]string
	names  []string
}

var enums = sync.OnceValue(func() enumData {
	data := enumData{values: make(map[string][]string)}

	gjson.GetBytes(openapi, "components.schemas").ForEach(func(key, value gjson.Result) bool {
		if value.Get("type").String() != "string" {
			return true
		}

		values := value.Get("enum")
		if !values.IsArray() {
			//+gocover:ignore:block OpenAPI always contain an enum array
			return true
		}

		name := key.String()
		result := make([]string, 0, len(values.Array()))
		for _, value := range values.Array() {
			result = append(result, value.String())
		}

		data.values[name] = result
		data.names = append(data.names, name)
		return true
	})

	return data
})

func EnumNames() []string {
	return enums().names
}

func Enum(name string) []string {
	return enums().values[name]
}
