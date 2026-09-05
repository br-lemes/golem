package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionFight(name string, participants []string) (schemas.CharacterFightDataSchema, error) {
	release := beginCriticalAction()
	defer release()
	resp, err := post(fmt.Sprintf("/my/%s/action/fight", name), schemas.FightRequestSchema{
		Participants: &participants,
	})
	if err != nil {
		return schemas.CharacterFightDataSchema{}, err
	}
	var data schemas.CharacterFightResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterFightDataSchema{}, err
	}
	cache.SaveCharacters(data.Data.Characters)
	release()
	if data.Data.Fight.Result == "win" {
		for _, character := range data.Data.Fight.Characters {
			if len(data.Data.Fight.Characters) > 1 {
				console.Printf("[%s] ", character.CharacterName)
			}
			console.Printf("XP gained: %d", character.Xp)
			if character.Gold > 0 {
				console.Printf(", Gold gained: %d", character.Gold)
			}
			printDropSchema(character.Drops)
			console.Printf("\n")
		}
	} else {
		console.Printf("💀 Fight lost!\n")
	}
	handleCooldown(data.Data.Cooldown.TotalSeconds, string(data.Data.Cooldown.Reason))
	return data.Data, nil
}
