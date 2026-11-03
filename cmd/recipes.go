package cmd

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type recipesFlags struct {
	Details bool     `flag:"details" desc:"Include direct recipe ingredients"`
	Without []string `flag:"without" desc:"Exclude an ingredient and recipes that depend on it"`
}

var recipesCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "recipes <skill>",
	Short: "List recipes with the highest crafting XP",
	Long: `List recipes with the highest crafting XP

Arguments:
  skill   The crafting skill to inspect.`,
	ValidArgsFunction: completion.Custom(1, func() []string { return cheapestSkills }).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !slices.Contains(cheapestSkills, args[0]) {
			return fmt.Errorf("invalid skill %q: allowed values are %v", args[0], cheapestSkills)
		}
		flags, err := utils.ReadFlags[recipesFlags](cmd)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return recipesRun(args[0], flags)
	},
}

func recipesRun(skill string, flags recipesFlags) error {
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return err
	}
	maxLevel := 0
	for _, character := range characters {
		level, _ := utils.GetCharacterCraftingSkillLevel(character, skill)
		maxLevel = max(maxLevel, level)
	}

	items := recipesItems(skill, maxLevel, catalog.Items().All(), flags.Without)
	if flags.Details {
		return console.Auto(recipeDetails(items))
	}
	levels := make([]map[string]int, 0, len(items))
	for _, item := range items {
		levels = append(levels, map[string]int{item.Code: *item.Craft.Level})
	}
	return console.Auto(levels)
}

func recipesItems(skill string, playerLevel int, items []*schemas.ItemSchema, without []string) []*schemas.ItemSchema {
	minLevel := max(1, playerLevel-10)
	excluded := make(map[string]bool, len(without))
	itemByCode := make(map[string]*schemas.ItemSchema, len(items))
	for _, code := range without {
		excluded[code] = true
	}
	for _, item := range items {
		itemByCode[item.Code] = item
	}

	result := make([]*schemas.ItemSchema, 0)
	for _, item := range items {
		if !isCraftable(*item) || *item.Craft.Skill != skill {
			continue
		}
		if *item.Craft.Level < minLevel || *item.Craft.Level > playerLevel {
			continue
		}
		if dependsOnExcluded(item.Code, excluded, itemByCode, map[string]bool{}) {
			continue
		}
		result = append(result, item)
	}
	slices.SortFunc(result, func(a, b *schemas.ItemSchema) int {
		xpA := CalculateArtifactsXP(*a.Craft.Level, playerLevel, skill, 0, *a.Craft.Items)
		xpB := CalculateArtifactsXP(*b.Craft.Level, playerLevel, skill, 0, *b.Craft.Items)
		if xpA != xpB {
			return cmp.Compare(xpB, xpA)
		}
		if *a.Craft.Level != *b.Craft.Level {
			return cmp.Compare(*b.Craft.Level, *a.Craft.Level)
		}
		return cmp.Compare(a.Code, b.Code)
	})
	return result
}

func dependsOnExcluded(code string, excluded map[string]bool, itemByCode map[string]*schemas.ItemSchema, visiting map[string]bool) bool {
	if excluded[code] {
		return true
	}
	if visiting[code] {
		return false
	}
	visiting[code] = true
	defer delete(visiting, code)

	item, found := itemByCode[code]
	if !found {
		item, found = catalog.Items().Get(code)
	}
	if !found || item.Craft == nil || item.Craft.Items == nil {
		return false
	}
	for _, ingredient := range *item.Craft.Items {
		if excluded[ingredient.Code] || dependsOnExcluded(ingredient.Code, excluded, itemByCode, visiting) {
			return true
		}
	}
	return false
}

func recipeDetails(items []*schemas.ItemSchema) []map[string]map[string]int {
	result := make([]map[string]map[string]int, 0, len(items))
	for _, item := range items {
		ingredients := make(map[string]int, len(*item.Craft.Items))
		for _, ingredient := range *item.Craft.Items {
			ingredients[ingredient.Code] = ingredient.Quantity
		}
		details := map[string]map[string]int{item.Code: ingredients}
		result = append(result, details)
	}
	return result
}

func init() {
	rootCmd.AddCommand(recipesCmd)
	err := utils.RegisterFlags[recipesFlags](recipesCmd)
	if err != nil {
		panic(err)
	}
	err = recipesCmd.RegisterFlagCompletionFunc("without", completion.UniqueList(catalog.Items().Keys))
	if err != nil {
		panic(err)
	}
}
