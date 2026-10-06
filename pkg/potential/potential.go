package potential

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

type Result struct {
	Code              string              `json:"code"`
	Used              bool                `json:"used"`
	Uses              []Use               `json:"uses,omitempty"`
	SelectedEquipment []SelectedEquipment `json:"selected_equipment,omitempty"`
}

type Use struct {
	CharacterLevel int                          `json:"character_level"`
	EquipmentLevel int                          `json:"equipment_level"`
	Monsters       []string                     `json:"monsters,omitempty"`
	Effects        []string                     `json:"effects,omitempty"`
	Loadouts       map[string]map[string]string `json:"loadouts,omitempty"`
}

type SelectedEquipment struct {
	Code string `json:"code"`
	Uses []Use  `json:"uses"`
}

type scenario struct {
	characterLevel int
	gearLevel      int
}

func Evaluate(codes, without []string, includeLoadouts bool, characterLevel, equipmentLevel int) (map[string]Result, error) {
	return evaluate(defaultDeps, codes, without, includeLoadouts, characterLevel, equipmentLevel)
}

func evaluate(d deps, codes, without []string, includeLoadouts bool, characterLevel, equipmentLevel int) (map[string]Result, error) {
	if characterLevel < 0 || equipmentLevel < 0 {
		return nil, fmt.Errorf("levels must not be negative")
	}
	cache.BeginFightSimulationBatch()
	defer cache.FlushFightSimulationBatch()
	targets := map[string]schemas.ItemSchema{}
	for _, code := range codes {
		item, ok := catalog.Items().Get(code)
		if !ok || !isEquipment(*item) {
			return nil, fmt.Errorf("item is not equipment: %s", code)
		}
		targets[code] = *item
	}
	excluded := make(map[string]bool, len(without))
	for _, excludedCode := range without {
		_, evaluated := targets[excludedCode]
		if evaluated {
			return nil, fmt.Errorf("cannot exclude evaluated item: %s", excludedCode)
		}
		excluded[excludedCode] = true
	}
	characters, err := d.accountsCharacters("")
	if err != nil {
		return nil, err
	}
	if characterLevel > 0 {
		for index := range characters {
			characters[index].Level = min(characters[index].Level, characterLevel)
		}
	}
	for code, item := range targets {
		canEquip := false
		for _, character := range characters {
			if canEquipAtLevel(character, item) {
				canEquip = true
			}
		}
		if !canEquip {
			return nil, fmt.Errorf("no character can equip item: %s", code)
		}
	}
	results := map[string]*Result{}
	selections := map[string]map[string]map[scenario]Use{}
	for code := range targets {
		results[code] = &Result{Code: code}
		selections[code] = map[string]map[scenario]Use{}
	}
	monsters := catalog.Monsters.Filter(func(monster *schemas.MonsterSchema) bool {
		return monster.Type != "boss" && monster.Type != "raid_boss"
	})
	slices.SortFunc(monsters, func(i, j *schemas.MonsterSchema) int {
		return cmp.Compare(i.Code, j.Code)
	})
	for _, current := range scenarios(characters, equipmentLevel) {
		character := schemas.CharacterSchema{Level: current.characterLevel}
		available := availableEquipment(character, excluded, current.gearLevel, targets)
		loadouts, findErr := cachedCombat(d, character, monsters, available)
		if findErr != nil {
			return nil, findErr
		}
		for _, item := range targets {
			if !canEquipAtLevel(character, item) {
				continue
			}
			console.Debugf("potential: character_level=%d gear_level=%d %s selections: %s\n", current.characterLevel, current.gearLevel, item.Type, selectedByType(monsters, loadouts, item.Type))
		}
		for _, monster := range monsters {
			selected, won := loadouts[monster.Code]
			if !won {
				continue
			}
			for code, item := range targets {
				if !canEquipAtLevel(character, item) {
					continue
				}
				if loadoutContains(selected, code) {
					use := resultUse(results[code], current)
					use.Monsters = appendUnique(use.Monsters, monster.Code)
					if includeLoadouts {
						if use.Loadouts == nil {
							use.Loadouts = map[string]map[string]string{}
						}
						use.Loadouts[monster.Code] = selected
					}
					setResultUse(results[code], use)
				}
				addSelections(selections[code], selected, item.Type, code, current, monster.Code, "", includeLoadouts)
			}
		}
		for _, priority := range []string{"wisdom", "prospecting"} {
			selected, findErr := d.findEquipment(character, best.EquipmentOptions{
				UniqueAdeptRing: true,
				Owned:           available,
				Priorities:      []string{priority},
			})
			if findErr != nil {
				return nil, fmt.Errorf("check %s at character level %d: %w", priority, current.characterLevel, findErr)
			}
			for code, item := range targets {
				if !canEquipAtLevel(character, item) {
					continue
				}
				if bestResultContains(selected, code) {
					use := resultUse(results[code], current)
					use.Effects = appendUnique(use.Effects, priority)
					setResultUse(results[code], use)
				}
				for _, selectedItem := range selected {
					addSelection(selections[code], selectedItem.Code, item.Type, code, current, "", priority, nil)
				}
			}
		}
	}
	final := map[string]Result{}
	for code, result := range results {
		result.Used = len(result.Uses) > 0
		if result.Used {
			final[code] = *result
			continue
		}
		for selectedCode, uses := range selections[code] {
			selected := SelectedEquipment{Code: selectedCode}
			for _, use := range uses {
				selected.Uses = append(selected.Uses, use)
			}
			slices.SortFunc(selected.Uses, func(i, j Use) int {
				return cmp.Compare(i.CharacterLevel, j.CharacterLevel)
			})
			result.SelectedEquipment = append(result.SelectedEquipment, selected)
		}
		slices.SortFunc(result.SelectedEquipment, func(i, j SelectedEquipment) int {
			return cmp.Compare(i.Code, j.Code)
		})
		final[code] = *result
	}
	return final, nil
}

func cachedCombat(d deps, character schemas.CharacterSchema, monsters []*schemas.MonsterSchema, available map[string]int) (map[string]map[string]string, error) {
	key := combatCacheKey(character.Level, monsters, available)
	stored, ok := cache.GetUsageCombat(cache.PotentialCombatType, key)
	if ok {
		var loadouts map[string]map[string]string
		err := json.Unmarshal([]byte(stored.Results), &loadouts)
		if err == nil {
			console.Debugf("potential: combat cache hit level=%d\n", character.Level)
			return loadouts, nil
		}
	}
	console.Debugf("potential: combat cache miss level=%d\n", character.Level)
	loadouts := map[string]map[string]string{}
	for _, monster := range monsters {
		fight, err := d.findFight(character, *monster, available, false, false)
		if err != nil {
			return nil, fmt.Errorf("check monster %s at level %d: %w", monster.Code, character.Level, err)
		}
		if fight.Winrate == 100 {
			loadouts[monster.Code] = fight.FinalEquipment
		}
	}
	encoded, err := json.Marshal(loadouts)
	if err != nil {
		//+gocover:ignore:block loadouts only contain JSON-compatible values
		return nil, err
	}
	cache.SaveUsageCombat(models.UsageCombat{
		Type:    cache.PotentialCombatType,
		Key:     key,
		Version: cache.UsageCombatVersion,
		Results: string(encoded),
	})
	return loadouts, nil
}

func combatCacheKey(level int, monsters []*schemas.MonsterSchema, available map[string]int) string {
	monsterCodes := make([]string, 0, len(monsters))
	for _, monster := range monsters {
		monsterCodes = append(monsterCodes, monster.Code)
	}
	codes := make([]string, 0, len(available))
	for code, quantity := range available {
		if quantity > 0 {
			codes = append(codes, code+"="+strconv.Itoa(quantity))
		}
	}
	slices.Sort(monsterCodes)
	slices.Sort(codes)
	input := strconv.Itoa(level) + "|" + strings.Join(monsterCodes, ",") + "|" + strings.Join(codes, ",")
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func selectedByType(monsters []*schemas.MonsterSchema, loadouts map[string]map[string]string, itemType string) string {
	selected := map[string][]string{}
	for _, monster := range monsters {
		loadout, ok := loadouts[monster.Code]
		if !ok {
			continue
		}
		for _, code := range loadout {
			item, ok := catalog.Items().Get(code)
			if ok && item.Type == itemType {
				selected[code] = appendUnique(selected[code], monster.Code)
			}
		}
	}
	codes := slices.Collect(maps.Keys(selected))
	slices.Sort(codes)
	parts := make([]string, 0, len(codes))
	for _, code := range codes {
		parts = append(parts, code+"="+strings.Join(selected[code], ","))
	}
	return strings.Join(parts, " ")
}

func scenarios(characters []schemas.CharacterSchema, gearLimit int) []scenario {
	set := map[scenario]bool{}
	for _, character := range characters {
		gearLevel := character.Level
		if gearLimit > 0 {
			gearLevel = min(gearLevel, gearLimit)
		}
		set[scenario{character.Level, gearLevel}] = true
	}
	result := slices.Collect(maps.Keys(set))
	slices.SortFunc(result, func(i, j scenario) int {
		return cmp.Compare(i.characterLevel, j.characterLevel)
	})
	return result
}

func availableEquipment(character schemas.CharacterSchema, excluded map[string]bool, gearLevel int, evaluated map[string]schemas.ItemSchema) map[string]int {
	available := map[string]int{}
	for _, item := range catalog.Items().All() {
		if excluded[item.Code] || !isAvailable(*item) || !canEquipAtLevel(character, *item) {
			continue
		}
		_, isEvaluated := evaluated[item.Code]
		if !isEvaluated && item.Level > gearLevel {
			continue
		}
		available[item.Code] = quantityLimit(item.Code, *item)
	}
	return available
}

func canEquipAtLevel(character schemas.CharacterSchema, item schemas.ItemSchema) bool {
	return item.Level <= character.Level && best.CanEquip(character, item)
}

func resultUse(result *Result, current scenario) Use {
	for _, use := range result.Uses {
		if use.CharacterLevel == current.characterLevel && use.EquipmentLevel == current.gearLevel {
			return use
		}
	}
	return Use{
		CharacterLevel: current.characterLevel,
		EquipmentLevel: current.gearLevel,
	}
}

func setResultUse(result *Result, use Use) {
	for index, current := range result.Uses {
		if current.CharacterLevel == use.CharacterLevel && current.EquipmentLevel == use.EquipmentLevel {
			result.Uses[index] = use
			return
		}
	}
	result.Uses = append(result.Uses, use)
}

func isAvailable(item schemas.ItemSchema) bool {
	return isEquipment(item) || item.Type == "utility"
}

func isEquipment(item schemas.ItemSchema) bool {
	_, ok := catalog.EquipmentTypeToSlots[item.Type]
	return ok && item.Type != "utility" && item.Subtype != "tool"
}

func quantityLimit(code string, item schemas.ItemSchema) int {
	if code == "ring_of_the_adept" {
		return 6
	}
	if item.Type == "ring" {
		return 10
	}
	return 5
}

func loadoutContains(loadout map[string]string, code string) bool {
	for _, selected := range loadout {
		if selected == code {
			return true
		}
	}
	return false
}

func bestResultContains(results map[string]best.BestResult, code string) bool {
	for _, selected := range results {
		if selected.Code == code {
			return true
		}
	}
	return false
}

func addSelections(selections map[string]map[scenario]Use, loadout map[string]string, itemType, evaluated string, current scenario, monster, effect string, details bool) {
	for _, code := range loadout {
		addSelection(selections, code, itemType, evaluated, current, monster, effect, loadoutForDetails(loadout, details))
	}
}

func addSelection(selections map[string]map[scenario]Use, code, itemType, evaluated string, current scenario, monster, effect string, loadout map[string]string) {
	if code == evaluated {
		return
	}
	item, ok := catalog.Items().Get(code)
	if ok && item.Type == itemType {
		if selections[code] == nil {
			selections[code] = map[scenario]Use{}
		}
		use := selections[code][current]
		if use.CharacterLevel == 0 {
			use = Use{
				CharacterLevel: current.characterLevel,
				EquipmentLevel: current.gearLevel,
			}
		}
		if monster != "" {
			use.Monsters = appendUnique(use.Monsters, monster)
			if loadout != nil {
				if use.Loadouts == nil {
					use.Loadouts = map[string]map[string]string{}
				}
				use.Loadouts[monster] = loadout
			}
		}
		if effect != "" {
			use.Effects = appendUnique(use.Effects, effect)
		}
		selections[code][current] = use
	}
}

func loadoutForDetails(loadout map[string]string, details bool) map[string]string {
	if details {
		return loadout
	}
	return nil
}

func appendUnique(values []string, value string) []string {
	if slices.Contains(values, value) {
		return values
	}
	return append(values, value)
}
