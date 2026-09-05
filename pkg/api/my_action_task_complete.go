package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionTaskComplete(name string) (schemas.RewardDataSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/task/complete", name)
	resp, err := post(path, nil)
	if err != nil {
		return schemas.RewardDataSchema{}, err
	}
	var data schemas.RewardDataResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.RewardDataSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
