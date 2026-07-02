package api

import (
	"encoding/json"
	"fmt"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

func MyActionFight(name string, participants []string) (schemas.CharacterFightDataSchema, error) {
	resp, err := PostNoCooldown(fmt.Sprintf("/my/%s/action/fight", name),
		schemas.FightRequestSchema{Participants: &participants})
	if err != nil {
		return schemas.CharacterFightDataSchema{}, err
	}
	var data schemas.CharacterFightResponseSchema
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return schemas.CharacterFightDataSchema{}, err
	}
	for _, character := range data.Data.Characters {
		cache.SaveCharacter(character.Name, character)
	}
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
		}
	} else {
		console.Printf("💀 Fight lost!")
	}
	console.Printf("\n")
	handleCooldown(data.Data.Cooldown.TotalSeconds)
	return data.Data, nil
}
