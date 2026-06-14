package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionRecycling(name string, item schemas.SimpleItemSchema) (schemas.RecyclingDataSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/recycling", name), item)
	if err != nil {
		return schemas.RecyclingDataSchema{}, err
	}
	var data schemas.RecyclingResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.RecyclingDataSchema{}, err
	}
	return data.Data, nil
}
