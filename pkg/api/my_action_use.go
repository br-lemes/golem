package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionUse(name string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/use", name), item)
	if err != nil {
		return schemas.UseItemSchema{}, err
	}
	var data schemas.UseItemResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.UseItemSchema{}, err
	}
	return data.Data, nil
}
