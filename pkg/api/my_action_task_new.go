package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionTaskNew(name string) (schemas.TaskDataSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/task/new", name)
	resp, err := PostNoCooldown(path, nil)
	if err != nil {
		return schemas.TaskDataSchema{}, err
	}
	var data schemas.TaskResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.TaskDataSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
