package datatable

import (
	"cmp"
	"net/http"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Page[T any] struct {
	Data       []T
	Page       int
	TotalItems int
}

func QueryFromRequest(request *http.Request) Query {
	values := request.URL.Query()
	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		page = 1
	}
	size := DefaultSize
	parsedSize, parseErr := strconv.Atoi(values.Get("size"))
	if parseErr == nil && slices.Contains(Sizes, parsedSize) {
		size = parsedSize
	}
	filters := make(map[string]string)
	for key, value := range values {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") && len(value) > 0 {
			filters[strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]")] = value[0]
		}
	}
	return Query{
		Page:      page,
		Size:      size,
		Search:    values.Get("search"),
		Sort:      values.Get("sort"),
		Direction: values.Get("direction"),
		Filters:   filters,
	}
}

func Search[T any](items []T, query string, columns []Column[T]) []T {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return items
	}
	filtered := make([]T, 0, len(items))
	for _, item := range items {
		value := reflect.ValueOf(item)
		for _, column := range columns {
			field := fieldValue(value, column.Field)
			if fieldSearchText(field) != "" && strings.Contains(fieldSearchText(field), query) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}

func Sort[T any](items []T, field, direction string, columns []Column[T]) {
	if direction != "asc" && direction != "desc" {
		return
	}
	for _, column := range columns {
		if column.Field != field || !column.Sortable {
			continue
		}
		sort.SliceStable(items, func(left, right int) bool {
			comparison := compareFields(reflect.ValueOf(items[left]), reflect.ValueOf(items[right]), field)
			if direction == "asc" {
				return comparison < 0
			}
			return comparison > 0
		})
		return
	}
}

func Paginate[T any](items []T, page, size int) ([]T, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = DefaultSize
	}
	totalPages := max(1, (len(items)+size-1)/size)
	page = min(page, totalPages)
	start := (page - 1) * size
	end := min(start+size, len(items))
	return items[start:end], page
}

func Prepare[T any](items []T, query Query, columns []Column[T]) Page[T] {
	Sort(items, query.Sort, query.Direction, columns)
	totalItems := len(items)
	items, page := Paginate(items, query.Page, query.Size)
	return Page[T]{Data: items, Page: page, TotalItems: totalItems}
}

func fieldSearchText(value reflect.Value) string {
	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return ""
	}
	var text string
	switch value.Kind() {
	case reflect.String:
		text = value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		text = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		text = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		text = strconv.FormatFloat(value.Float(), 'f', -1, value.Type().Bits())
	default:
		return ""
	}
	return strings.ToLower(text)
}

func compareFields(left, right reflect.Value, fieldName string) int {
	left = fieldValue(left, fieldName)
	right = fieldValue(right, fieldName)
	if !left.IsValid() || !right.IsValid() {
		return 0
	}
	switch left.Kind() {
	case reflect.String:
		return strings.Compare(left.String(), right.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cmp.Compare(left.Int(), right.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cmp.Compare(left.Uint(), right.Uint())
	case reflect.Float32, reflect.Float64:
		return cmp.Compare(left.Float(), right.Float())
	case reflect.Bool:
		if left.Bool() == right.Bool() {
			return 0
		}
		if !left.Bool() {
			return -1
		}
		return 1
	default:
		return 0
	}
}

func fieldValue(value reflect.Value, fieldName string) reflect.Value {
	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	if !value.IsValid() || value.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	typeOfValue := value.Type()
	for index := 0; index < value.NumField(); index++ {
		field := typeOfValue.Field(index)
		if strings.Split(field.Tag.Get("json"), ",")[0] == fieldName {
			return value.Field(index)
		}
	}
	return reflect.Value{}
}
