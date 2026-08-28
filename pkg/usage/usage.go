package usage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

// CACHE_VERSION / ACTION REQUIRED: increment when usage logic, its cached
// result format, or a dependency changes the selected combat loadouts.
// This cache also depends on bestFightSimulationCacheVersion; increment this
// version too whenever that simulator version is incremented.
const combatCacheVersion = 1

type Evaluation struct {
	Needed   bool     `json:"needed"`
	Monsters []string `json:"monsters,omitempty"`
	Effects  []string `json:"effects,omitempty"`
}

func Evaluate(codes []string, details bool) (map[string]Evaluation, error) {
	if len(codes) > 0 {
		seen := map[string]bool{}
		for _, code := range codes {
			if seen[code] {
				continue
			}
			item, ok := database.Items().Get(code)
			if !ok || !isEquipment(*item) {
				return nil, fmt.Errorf("item is not equipment: %s", code)
			}
			seen[code] = true
		}
	}
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return nil, err
	}
	if len(characters) == 0 {
		return nil, fmt.Errorf("account has no characters")
	}
	character := characters[0]
	for _, candidate := range characters[1:] {
		if candidate.Level < character.Level {
			character = candidate
		}
	}
	simulationCharacter := schemas.CharacterSchema{Level: character.Level}
	owned, err := globalOwned(characters)
	if err != nil {
		return nil, err
	}
	allCodes := equipmentCodes(simulationCharacter, owned)
	if len(codes) == 0 {
		codes = allCodes
	} else {
		for _, code := range codes {
			if owned[code] < 1 {
				return nil, fmt.Errorf("item is not owned: %s", code)
			}
		}
	}
	combatAvailable := combatItems(simulationCharacter, owned)
	result := make(map[string]Evaluation, len(codes))
	for _, code := range codes {
		item, ok := database.Items().Get(code)
		if !ok || !isEquipment(*item) || !best.CanEquip(simulationCharacter, *item) {
			continue
		}
		result[code] = Evaluation{}
	}

	monsters := database.Monsters.Filter(func(m *schemas.MonsterSchema) bool { return m.Type != schemas.Boss && m.Type != schemas.RaidBoss })
	sort.Slice(monsters, func(i, j int) bool { return monsters[i].Code < monsters[j].Code })
	normalCombat, err := cachedMarkCombat(simulationCharacter, monsters, combatAvailable, "normal")
	if err != nil {
		return nil, err
	}
	mergeCombatUsage(result, normalCombat, details)
	err = markCrafting(simulationCharacter, owned, result, details)
	if err != nil {
		return nil, err
	}

	for _, code := range insufficientCodes(simulationCharacter, owned) {
		available := cloneCounts(combatAvailable)
		delete(available, code)
		shortageCombat, err := cachedMarkCombat(simulationCharacter, monsters, available, "shortage of "+code)
		if err != nil {
			return nil, err
		}
		mergeCombatUsage(result, shortageCombat, details)
		err = markCrafting(simulationCharacter, available, result, details)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func cachedMarkCombat(character schemas.CharacterSchema, monsters []*schemas.MonsterSchema, available map[string]int, scenario string) (map[string]map[string]string, error) {
	key := combatCacheKey(character.Level, monsters, available)
	console.Debugf("usage: combat cache scenario=%s key=%s available=%s\n", scenario, key, canonicalAvailable(available))
	stored, ok := cache.GetUsageCombat(key, combatCacheVersion)
	if ok {
		var loadouts map[string]map[string]string
		err := json.Unmarshal([]byte(stored.Results), &loadouts)
		if err == nil {
			return loadouts, nil
		}
	}
	loadouts, err := markCombat(character, monsters, available, scenario)
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(loadouts)
	if err != nil {
		return nil, err
	}
	cache.SaveUsageCombat(models.UsageCombat{
		Key:     key,
		Version: combatCacheVersion,
		Results: string(encoded),
	})
	return loadouts, nil
}

func combatCacheKey(level int, monsters []*schemas.MonsterSchema, available map[string]int) string {
	monsterCodes := make([]string, 0, len(monsters))
	for _, monster := range monsters {
		monsterCodes = append(monsterCodes, monster.Code)
	}
	sort.Strings(monsterCodes)
	input := strconv.Itoa(level) + "|" + strings.Join(monsterCodes, ",") + "|" + canonicalAvailable(available)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func canonicalAvailable(available map[string]int) string {
	codes := make([]string, 0, len(available))
	for code, quantity := range available {
		if quantity > 0 {
			canonicalQuantity := quantity
			item, ok := database.Items().Get(code)
			if ok {
				if item.Type == "utility" {
					canonicalQuantity = 1
				} else if isEquipment(*item) {
					canonicalQuantity = min(canonicalQuantity, equipmentQuantityLimit(code, *item))
				}
			}
			codes = append(codes, code+"="+strconv.Itoa(canonicalQuantity))
		}
	}
	sort.Strings(codes)
	return strings.Join(codes, ",")
}

func equipmentQuantityLimit(code string, item schemas.ItemSchema) int {
	if code == "ring_of_the_adept" {
		return 6
	}
	if item.Type == "ring" {
		return 10
	}
	return 5
}

func markCombat(character schemas.CharacterSchema, monsters []*schemas.MonsterSchema, available map[string]int, scenario string) (map[string]map[string]string, error) {
	loadouts := map[string]map[string]string{}
	for _, monster := range monsters {
		console.Debugf("usage: checking %s against monster %s\n", scenario, monster.Code)
		fight, err := best.FindFightWithAvailable(character, *monster, available, false, false)
		if err != nil {
			return nil, fmt.Errorf("check monster %s (%s): %w", monster.Code, scenario, err)
		}
		if fight.Winrate < 100 {
			continue
		}
		loadouts[monster.Code] = fight.FinalEquipment
	}
	return loadouts, nil
}

func mergeCombatUsage(result map[string]Evaluation, loadouts map[string]map[string]string, details bool) {
	for monster, loadout := range loadouts {
		for code := range result {
			if !loadoutContains(loadout, code) {
				continue
			}
			current := result[code]
			current.Needed = true
			if details && !contains(current.Monsters, monster) {
				current.Monsters = append(current.Monsters, monster)
			}
			result[code] = current
		}
	}
}

func markCrafting(character schemas.CharacterSchema, available map[string]int, result map[string]Evaluation, details bool) error {
	for _, priority := range []string{"wisdom", "prospecting"} {
		selected, err := best.FindEquipment(character, best.EquipmentOptions{
			UniqueAdeptRing: true,
			Owned:           available,
			Priorities:      []string{priority},
		})
		if err != nil {
			return fmt.Errorf("check crafting priority %s: %w", priority, err)
		}
		loadout := characterLoadout(character)
		for slot, item := range selected {
			loadout[slot] = item.Code
		}
		for code := range result {
			if loadoutContains(loadout, code) {
				current := result[code]
				current.Needed = true
				if details && !contains(current.Effects, priority) {
					current.Effects = append(current.Effects, priority)
				}
				result[code] = current
			}
		}
	}
	return nil
}

func globalOwned(characters []schemas.CharacterSchema) (map[string]int, error) {
	owned := map[string]int{}
	bank, err := api.MyBankItems()
	if err != nil {
		return nil, err
	}
	for _, item := range bank {
		owned[item.Code] += item.Quantity
	}
	for _, character := range characters {
		for _, code := range characterLoadout(character) {
			if code != "" {
				owned[code]++
			}
		}
		if character.Inventory != nil {
			for _, item := range *character.Inventory {
				owned[item.Code] += item.Quantity
			}
		}
	}
	return owned, nil
}

func equipmentCodes(character schemas.CharacterSchema, owned map[string]int) []string {
	codes := []string{}
	for code := range owned {
		item, ok := database.Items().Get(code)
		if ok && isEquipment(*item) && best.CanEquip(character, *item) {
			codes = append(codes, code)
		}
	}
	sort.Strings(codes)
	return codes
}

func insufficientCodes(character schemas.CharacterSchema, owned map[string]int) []string {
	codes := []string{}
	for _, code := range equipmentCodes(character, owned) {
		item, _ := database.Items().Get(code)
		required := 5
		if item.Type == "ring" {
			required = 10
			if code == "ring_of_the_adept" {
				required = 5
			}
		}
		if owned[code] < required {
			codes = append(codes, code)
		}
	}
	return codes
}

func isEquipment(item schemas.ItemSchema) bool {
	_, ok := database.EquipmentTypeToSlots[item.Type]
	return ok && item.Type != "utility" && item.Subtype != "tool"
}

func combatItems(character schemas.CharacterSchema, owned map[string]int) map[string]int {
	available := map[string]int{}
	for code, quantity := range owned {
		item, ok := database.Items().Get(code)
		if !ok || quantity < 1 {
			continue
		}
		if (isEquipment(*item) || item.Type == "utility") && best.CanEquip(character, *item) {
			available[code] = quantity
		}
	}
	return available
}
func cloneCounts(source map[string]int) map[string]int {
	result := map[string]int{}
	for code, quantity := range source {
		result[code] = quantity
	}
	return result
}
func loadoutContains(loadout map[string]string, code string) bool {
	for _, selected := range loadout {
		if selected == code {
			return true
		}
	}
	return false
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
func characterLoadout(c schemas.CharacterSchema) map[string]string {
	return map[string]string{
		"amulet":     c.AmuletSlot,
		"artifact1":  c.Artifact1Slot,
		"artifact2":  c.Artifact2Slot,
		"artifact3":  c.Artifact3Slot,
		"bag":        c.BagSlot,
		"body_armor": c.BodyArmorSlot,
		"boots":      c.BootsSlot,
		"helmet":     c.HelmetSlot,
		"leg_armor":  c.LegArmorSlot,
		"ring1":      c.Ring1Slot,
		"ring2":      c.Ring2Slot,
		"rune":       c.RuneSlot,
		"shield":     c.ShieldSlot,
		"utility1":   c.Utility1Slot,
		"utility2":   c.Utility2Slot,
		"weapon":     c.WeaponSlot,
	}
}
