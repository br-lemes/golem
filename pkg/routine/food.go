package routine

import (
	"sort"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

const foodMinThreshold = 5
const foodFillPercent = 70

func foodCurrentQty(character schemas.CharacterSchema) int {
	if character.Inventory == nil {
		return 0
	}
	total := 0
	for _, slot := range *character.Inventory {
		if slot.Quantity <= 0 {
			continue
		}
		item, found := database.Items.Get(slot.Code)
		if !found || item.Type != "consumable" || item.Subtype != "food" {
			continue
		}
		total += slot.Quantity
	}
	return total
}

func foodCheck(character schemas.CharacterSchema, food string, bankQty map[string]int) bool {
	if food == "" || foodCurrentQty(character) >= foodMinThreshold {
		return false
	}
	return len(foodCandidates(character, food, bankQty)) > 0
}

func foodCandidates(character schemas.CharacterSchema, food string, bankQty map[string]int) []string {
	type candidate struct {
		code string
		heal int
	}
	var candidates []candidate
	for code, qty := range bankQty {
		if qty <= 0 {
			continue
		}
		item, found := database.Items.Get(code)
		if !found || item.Type != "consumable" || item.Subtype != "food" || item.Level > character.Level || item.Effects == nil {
			continue
		}
		heal := 0
		for _, effect := range *item.Effects {
			if effect.Code == "heal" {
				heal = effect.Value
				break
			}
		}
		if heal <= 0 {
			continue
		}
		candidates = append(candidates, candidate{code, heal})
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].heal > candidates[j].heal })
	codes := make([]string, 0, len(candidates))
	if food != "" && food != "auto" {
		for _, c := range candidates {
			if c.code == food {
				codes = append(codes, food)
				break
			}
		}
	}
	for _, c := range candidates {
		if c.code == food {
			continue
		}
		codes = append(codes, c.code)
	}
	return codes
}

func SelectFood(character schemas.CharacterSchema, bankQty map[string]int) string {
	candidates := foodCandidates(character, "", bankQty)
	if len(candidates) == 0 {
		return ""
	}
	return candidates[0]
}

func foodRestock(d deps, character schemas.CharacterSchema, food string, bankQty map[string]int) (schemas.CharacterSchema, error) {
	if food == "" {
		return character, nil
	}
	codes := foodCandidates(character, food, bankQty)
	if len(codes) == 0 {
		console.Printf("  No suitable food available in bank\n")
		return character, nil
	}
	target := (character.InventoryMaxItems - totalItems(character)) * foodFillPercent / 100
	if target <= 0 {
		return character, nil
	}
	remaining, filled := target, 0
	for _, code := range codes {
		for remaining > 0 && bankQty[code] > 0 {
			space := character.InventoryMaxItems - totalItems(character)
			if space <= 0 {
				console.Printf("  No inventory space to withdraw more food\n")
				remaining = 0
				break
			}
			chunk := remaining
			if bankQty[code] < chunk {
				chunk = bankQty[code]
			}
			if space < chunk {
				chunk = space
			}
			withdrawData, err := d.myActionBankWithdrawItem(character.Name, []schemas.SimpleItemSchema{
				{Code: code, Quantity: chunk},
			})
			if err != nil {
				return character, err
			}
			character = withdrawData.Character
			bankQty[code] -= chunk
			remaining -= chunk
			filled += chunk
		}
		if remaining == 0 {
			break
		}
	}
	if filled < target {
		console.Printf("  Only %d/%d food stocked\n", filled, target)
	}
	return character, nil
}
