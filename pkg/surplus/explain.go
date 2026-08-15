package surplus

import (
	"fmt"
	"sort"

	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
)

type Explanation struct {
	Code        string
	Status      string
	Reason      string
	Item        schemas.ItemSchema
	Total       int
	Surplus     int
	DominatedBy []Comparison
	ComparedTo  []Comparison
}

type Comparison struct {
	Code    string
	Level   int
	Effects map[string]string
}

func Explain(input Input, code string) Explanation {
	owned := collectEquipment(input)
	item, ok := owned[code]
	if !ok {
		return Explanation{
			Code:   code,
			Status: "not_owned",
			Reason: "item is not currently owned",
		}
	}

	positions := itemPositions(item.Item, input.Characters)
	if len(positions) == 0 {
		return Explanation{
			Code:   code,
			Status: "future",
			Item:   item.Item,
			Total:  item.Total,
			Reason: "no character can use this item yet; it is being kept for future progression",
		}
	}

	for _, result := range Find(input) {
		if result.Item.Code == code {
			dominatedBy := make([]Comparison, 0, len(result.DominatedBy))
			for _, superiorCode := range result.DominatedBy {
				superior, ok := owned[superiorCode]
				if !ok {
					//+gocover:ignore:block impossible under current invariants
					continue
				}
				dominatedBy = append(dominatedBy, Comparison{
					Code:    superiorCode,
					Level:   superior.Item.Level,
					Effects: compareEffects(item.Item, superior.Item),
				})
			}
			reason := fmt.Sprintf("%d owned, %d needed", result.Total, result.Total-result.Surplus)
			if len(dominatedBy) > 0 {
				reason += "; some positions covered by other equipment"
			}
			return Explanation{
				Code:        code,
				Status:      "surplus",
				Item:        result.Item,
				Total:       result.Total,
				Surplus:     result.Surplus,
				DominatedBy: dominatedBy,
				Reason:      reason,
			}
		}
	}

	comparisons := []Comparison{}
	for superiorCode, superior := range owned {
		if superiorCode == code || !compatibleItems(superior.Item, item.Item) {
			continue
		}
		comparisons = append(comparisons, Comparison{
			Code:    superiorCode,
			Level:   superior.Item.Level,
			Effects: compareEffects(item.Item, superior.Item),
		})
	}
	sort.Slice(comparisons, func(i, j int) bool { return comparisons[i].Code < comparisons[j].Code })

	reason := "no owned compatible item dominates this item"
	if len(comparisons) == 0 {
		reason = "no other owned item has the same type and subtype"
	}
	return Explanation{
		Code:       code,
		Status:     "not_dominated",
		Item:       item.Item,
		Total:      item.Total,
		Reason:     reason,
		ComparedTo: comparisons,
	}
}

func Evaluate(input Input, code string) Explanation {
	_, ok := database.Items.Get(code)
	if !ok {
		return Explanation{
			Code:   code,
			Status: "unknown",
			Reason: "item is not in the catalog",
		}
	}
	simulated := input
	evaluated := schemas.SimpleItemSchema{Code: code, Quantity: 1}
	simulated.BankItems = append(append([]schemas.SimpleItemSchema{}, input.BankItems...), evaluated)
	result := Explain(simulated, code)
	return result
}

func compareEffects(item, other schemas.ItemSchema) map[string]string {
	itemEffects, otherEffects := effectValues(item), effectValues(other)
	codes := make(map[string]struct{}, len(itemEffects)+len(otherEffects))
	for code := range itemEffects {
		codes[code] = struct{}{}
	}
	for code := range otherEffects {
		codes[code] = struct{}{}
	}
	effects := make(map[string]string)
	for code := range codes {
		itemValue, otherValue := itemEffects[code], otherEffects[code]
		if itemValue == otherValue {
			continue
		}
		op := ">"
		if itemValue < otherValue {
			op = "<"
		}
		meaning := "worse"
		if atLeastAsGood(code, itemValue, otherValue) {
			meaning = "better"
		}
		effects[code] = fmt.Sprintf("%d %s %d (%s)", itemValue, op, otherValue, meaning)
	}
	return effects
}
