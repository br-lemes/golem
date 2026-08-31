package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionTaskTrade(name string, item schemas.SimpleItemSchema) (schemas.TaskTradeDataSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/task/trade", name)
	resp, err := PostNoCooldown(path, item)
	if err != nil {
		return schemas.TaskTradeDataSchema{}, err
	}
	var data schemas.TaskTradeResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.TaskTradeDataSchema{}, err
	}
	cache.SaveCharacter(data.Data.Character)
	release()
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
