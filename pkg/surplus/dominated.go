package surplus

import (
	"sort"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

func Find(input Input) []Result {
	owned := collectEquipment(input)
	for code, inferior := range owned {
		positions := itemPositions(inferior.Item, input.Characters)
		if len(positions) == 0 {
			positions = allItemPositions(inferior.Item, input.Characters)
		}

		type candidate struct {
			item      Result
			positions []int
		}
		candidates := []candidate{}
		for superiorCode, superior := range owned {
			if superiorCode == code {
				continue
			}
			replaceable := []int{}
			for i, characterIndex := range positions {
				if CanReplace(superior.Item, inferior.Item, input.Characters[characterIndex]) {
					replaceable = append(replaceable, i)
				}
			}
			if len(replaceable) > 0 {
				candidates = append(candidates, candidate{
					item:      superior,
					positions: replaceable,
				})
			}
		}
		sort.Slice(candidates, func(i, j int) bool {
			return len(candidates[i].positions) < len(candidates[j].positions)
		})

		available := make([]bool, len(positions))
		for i := range available {
			available[i] = true
		}
		covered := 0
		for _, candidate := range candidates {
			used := 0
			for _, index := range candidate.positions {
				if !available[index] || used >= candidate.item.Total {
					continue
				}
				available[index] = false
				used++
				covered++
			}
			if used > 0 {
				inferior.DominatedBy = append(inferior.DominatedBy, candidate.item.Item.Code)
			}
		}

		needed := len(positions) - covered
		surplus := inferior.Total - needed
		if surplus > 0 {
			inferior.Surplus = surplus
			owned[code] = inferior
		}
	}

	result := make([]Result, 0, len(owned))
	for _, item := range owned {
		if item.Surplus > 0 {
			result = append(result, item)
		}
	}
	return result
}

func allItemPositions(item schemas.ItemSchema, characters []schemas.CharacterSchema) []int {
	capacity := 1
	if item.Type == "ring" {
		capacity = 2
	}
	positions := make([]int, 0, len(characters)*capacity)
	for i := range characters {
		for range capacity {
			positions = append(positions, i)
		}
	}
	return positions
}

func itemPositions(item schemas.ItemSchema, characters []schemas.CharacterSchema) []int {
	capacity := 1
	if item.Type == "ring" {
		capacity = 2
	}

	positions := []int{}
	for i, character := range characters {
		if !CanUse(item, character) {
			continue
		}
		for range capacity {
			positions = append(positions, i)
		}
	}
	return positions
}

func collectEquipment(input Input) map[string]Result {
	itemsByCode := make(map[string]Result, len(database.Items.All()))
	for _, item := range database.Items.All() {
		if item.Type == "utility" {
			continue
		}
		_, isEquipment := database.EquipmentTypeToSlots[item.Type]
		if !isEquipment {
			continue
		}
		itemsByCode[item.Code] = Result{Item: *item}
	}

	add := func(code string, quantity int) {
		item, ok := itemsByCode[code]
		if !ok {
			return
		}
		item.Total += quantity
		itemsByCode[code] = item
	}
	for _, item := range input.BankItems {
		add(item.Code, item.Quantity)
	}
	for _, character := range input.Characters {
		if character.Inventory != nil {
			for _, item := range *character.Inventory {
				add(item.Code, item.Quantity)
			}
		}
		for code, quantity := range equippedItems(character) {
			add(code, quantity)
		}
	}

	result := map[string]Result{}
	for code, item := range itemsByCode {
		if item.Total > 0 {
			result[code] = item
		}
	}
	return result
}
