package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionCrafting(name string, item schemas.SimpleItemSchema) (schemas.SkillDataSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/crafting", name), item)
	if err != nil {
		return schemas.SkillDataSchema{}, err
	}
	var data schemas.SkillResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.SkillDataSchema{}, err
	}
	return data.Data, nil
}
