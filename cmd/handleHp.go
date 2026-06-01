package cmd

import (
	"fmt"
	"time"
)

func handleHp(character CharacterSchema) (CharacterSchema, error) {
	threshold := character.MaxHp - ((character.MaxHp * 3) / 100)
	fmt.Fprintf(writer, "[%s] HP: %d/%d (threshold: %d)\n",
		time.Now().Format("15:04:05"), character.Hp, character.MaxHp, threshold)
	if character.Hp > threshold {
		return character, nil
	}
	fmt.Fprintf(writer, "[%s] Starting rest...\n",
		time.Now().Format("15:04:05"))
	data, err := apiActionRest(character.Name)
	if err != nil {
		return CharacterSchema{}, err
	}
	return data.Character, nil
}
