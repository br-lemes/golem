package utils

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestSizes(t *testing.T) {
	size, err := GetSize("/items")
	if err != nil || size <= 0 {
		t.Fatalf("GetSize(items) = %d, %v", size, err)
	}
	_, err = GetSize("/does-not-exist")
	if err == nil {
		t.Fatal("GetSize accepted an unknown path")
	}
	sizes, err := GetSizes()
	if err != nil || sizes["/items"] != size {
		t.Fatalf("GetSizes() = %v, %v", sizes, err)
	}
}

func TestGetSizeWithoutSizeParameter(t *testing.T) {
	_, err := GetSize("/my/details")
	if err == nil {
		t.Fatal("GetSize accepted a route without a size parameter")
	}
}

func TestSizeMax(t *testing.T) {
	max, ok := sizeMax(&openapi3.Operation{})
	if ok || max != 0 {
		t.Fatalf("sizeMax(empty) = %d, %t", max, ok)
	}
	limit := float64(25)
	op := &openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{Value: &openapi3.Parameter{Name: "other"}},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name: "size",
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{Max: &limit},
					},
				},
			},
		},
	}
	got, ok := sizeMax(op)
	if !ok || got != 25 {
		t.Fatalf("sizeMax() = %d, %t", got, ok)
	}
}
