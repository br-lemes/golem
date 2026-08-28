package cmd

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var cheapestSkills = []string{
	"gearcrafting",
	"jewelrycrafting",
	"weaponcrafting",
}

var cheapestCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "cheapest <skill>",
	Short: "List the cheapest items for leveling a skill",
	Long: `List the cheapest items for leveling a skill

Arguments:
  skill   The crafting skill to level.`,
	ValidArgsFunction: completion.StringSlice(func() []string { return cheapestSkills }),
	RunE: func(cmd *cobra.Command, args []string) error {
		skill := args[0]
		err := cheapestValidate(skill)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return cheapestRun(skill)
	},
}

func cheapestValidate(skill string) error {
	if !slices.Contains(cheapestSkills, skill) {
		return fmt.Errorf("invalid skill %q: allowed values are %v", skill, cheapestSkills)
	}
	return nil
}

func cheapestRun(skill string) error {
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return err
	}
	maxLevel := 0
	for _, character := range characters {
		level, _ := utils.GetCharacterCraftingSkillLevel(character, skill)
		maxLevel = max(maxLevel, level)
	}
	items := cheapestItems(skill, maxLevel, database.Items().All())
	codes := make([]string, 0, len(items))
	for _, item := range items {
		codes = append(codes, item.Code)
	}
	return console.Auto(codes)
}

type craftCost struct {
	basic, drops, rarity, monsterLevel, steps, units int
}

func craftCostLess(a, b craftCost) bool {
	pairs := [][2]int{
		{a.basic, b.basic},
		{a.drops, b.drops},
		{a.rarity, b.rarity},
		{a.monsterLevel, b.monsterLevel},
		{a.steps, b.steps},
		{a.units, b.units},
	}
	for _, pair := range pairs {
		if pair[0] != pair[1] {
			return pair[0] < pair[1]
		}
	}
	return false
}

func itemCraftCost(code string, visiting map[string]bool) craftCost {
	item, ok := database.Items().Get(code)
	if !ok || item.Craft == nil || item.Craft.Items == nil {
		return sourceCost(code)
	}
	if visiting[code] {
		return craftCost{basic: 1, rarity: 1000, units: 1}
	}
	visiting[code] = true
	result := craftCost{steps: 1}
	for _, ingredient := range *item.Craft.Items {
		cost := itemCraftCost(ingredient.Code, visiting)
		result.basic += cost.basic * ingredient.Quantity
		result.drops += cost.drops * ingredient.Quantity
		result.rarity += cost.rarity * ingredient.Quantity
		result.monsterLevel += cost.monsterLevel * ingredient.Quantity
		result.steps += cost.steps
		result.units += cost.units * ingredient.Quantity
	}
	delete(visiting, code)
	return result
}

func sourceCost(code string) craftCost {
	for _, npcItem := range database.NpcsItems.All() {
		if npcItem.Code != code || npcItem.BuyPrice == nil {
			continue
		}
		currency := sourceCost(npcItem.Currency)
		currency.basic *= *npcItem.BuyPrice
		currency.drops *= *npcItem.BuyPrice
		currency.rarity *= *npcItem.BuyPrice
		currency.monsterLevel *= *npcItem.BuyPrice
		currency.units *= *npcItem.BuyPrice
		return currency
	}
	best := craftCost{basic: 1, units: 1, rarity: 1}
	found := false
	for _, monster := range database.Monsters.All() {
		for _, drop := range monster.Drops {
			if drop.Code != code {
				continue
			}
			// Lower drop rates and stronger monsters are more expensive.
			bossMultiplier := 1
			switch monster.Type {
			case "elite":
				bossMultiplier = 5
			case "boss":
				bossMultiplier = 10
			case "raid_boss":
				bossMultiplier = 20
			}
			candidate := craftCost{
				basic:        bossMultiplier * 100,
				drops:        1,
				rarity:       bossMultiplier * 100000 / max(1, drop.Rate),
				monsterLevel: bossMultiplier * monster.Level,
				units:        1,
			}
			if !found || craftCostLess(candidate, best) {
				best, found = candidate, true
			}
		}
	}
	if found {
		return best
	}
	_, ok := database.Resources.Get(code)
	if ok {
		return craftCost{basic: 1, units: 1}
	}
	// Unknown sources are deliberately expensive rather than silently cheap.
	return craftCost{basic: 1, rarity: 1000000, units: 1}
}

func cheapestItems(skill string, maxLevel int, items []*schemas.ItemSchema) []*schemas.ItemSchema {
	minLevel := max(0, maxLevel-10)
	result := make([]*schemas.ItemSchema, 0)
	for _, item := range items {
		if item.Craft == nil || item.Craft.Skill == nil || string(*item.Craft.Skill) != skill || item.Craft.Level == nil {
			continue
		}
		if dependsOnTasksCoin(item.Code, map[string]bool{}) {
			continue
		}
		if *item.Craft.Level >= minLevel && *item.Craft.Level <= maxLevel {
			result = append(result, item)
		}
	}
	slices.SortFunc(result, func(a, b *schemas.ItemSchema) int {
		costA := itemCraftCost(a.Code, map[string]bool{})
		costB := itemCraftCost(b.Code, map[string]bool{})
		if craftCostLess(costA, costB) {
			return -1
		}
		if craftCostLess(costB, costA) {
			return 1
		}
		return cmp.Compare(a.Code, b.Code)
	})
	return result
}

func dependsOnTasksCoin(code string, visiting map[string]bool) bool {
	if code == "tasks_coin" || visiting[code] {
		return code == "tasks_coin"
	}
	visiting[code] = true
	defer delete(visiting, code)

	for _, npcItem := range database.NpcsItems.All() {
		if npcItem.Code == code && npcItem.Currency == "tasks_coin" {
			return true
		}
	}
	item, ok := database.Items().Get(code)
	if !ok || item.Craft == nil || item.Craft.Items == nil {
		return false
	}
	for _, ingredient := range *item.Craft.Items {
		if dependsOnTasksCoin(ingredient.Code, visiting) {
			return true
		}
	}
	return false
}

func init() { rootCmd.AddCommand(cheapestCmd) }
