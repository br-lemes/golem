package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionTaskCancel(name string) (schemas.TaskCancelledSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/task/cancel", name)
	resp, err := PostNoCooldown(path, nil)
	if err != nil {
		return schemas.TaskCancelledSchema{}, err
	}
	var data schemas.TaskCancelledResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.TaskCancelledSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
