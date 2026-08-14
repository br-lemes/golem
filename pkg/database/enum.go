package database

import (
	"sync"

	"github.com/tidwall/gjson"
)

var (
	enumCache sync.Map
	enumList  []string
)

func GetEnumNames() []string {
	initEnumsList()
	return enumList
}

func GetEnum(name string) []string {
	cached, exists := enumCache.Load(name)
	if exists {
		if cached == nil {
			return nil
		}
		return cached.([]string)
	}
	result := getEnum(name)
	actual, _ := enumCache.LoadOrStore(name, result)
	if actual == nil {
		return nil
	}
	return actual.([]string)
}

func getEnum(name string) []string {
	schema := gjson.GetBytes(openapi, "components.schemas."+name)
	if schema.Get("type").String() != "string" {
		return nil
	}
	enum := schema.Get("enum")
	if !enum.IsArray() {
		return nil
	}
	var result []string
	for _, v := range enum.Array() {
		result = append(result, v.String())
	}
	return result
}

var initEnumsList = sync.OnceFunc(func() {
	gjson.GetBytes(openapi, "components.schemas").ForEach(func(key, value gjson.Result) bool {
		if value.Get("type").String() != "string" {
			return true
		}
		if !value.Get("enum").IsArray() {
			return true
		}
		var result []string
		for _, v := range value.Get("enum").Array() {
			result = append(result, v.String())
		}
		name := key.String()
		enumCache.Store(name, result)
		enumList = append(enumList, key.String())
		return true
	})
})
