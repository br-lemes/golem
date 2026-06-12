package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionGathering(name string) (schemas.SkillDataSchema, error) {
	resp, err := Post(fmt.Sprintf("/my/%s/action/gathering", name), nil)
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
