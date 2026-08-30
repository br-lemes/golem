package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionGathering(name string) (schemas.SkillDataSchema, error) {
	release := beginCriticalAction()
	defer release()
	path := fmt.Sprintf("/my/%s/action/gathering", name)
	resp, err := PostNoCooldown(path, nil)
	if err != nil {
		return schemas.SkillDataSchema{}, err
	}
	var data schemas.SkillResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.SkillDataSchema{}, err
	}
	cache.SaveCharacter(name, data.Data.Character)
	release()
	console.Printf("XP gained: %d", data.Data.Details.Xp)
	printDropSchema(data.Data.Details.Items)
	console.Printf("\n")
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
