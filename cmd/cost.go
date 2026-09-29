package cmd

import (
	"fmt"
	"math"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type costFlags struct {
	Quantity int `flag:"quantity" shorthand:"q" desc:"Amount of items to craft (0 for maximum available)"`
}

var costCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "cost <code>",
	Short: "Calculate resources and XP for crafting an item",
	Long: `Calculate resources and XP for crafting an item

Arguments:
  code   The code of the item.`,
	ValidArgsFunction: completion.Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]

		flags, err := utils.ReadFlags[costFlags](cmd)
		if err != nil {
			return err
		}
		if flags.Quantity < 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
		item, found := catalog.Items().Get(code)
		if !found {
			return fmt.Errorf("item not found: %s", code)
		}
		cmd.SilenceUsage = true

		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}

		if !isCraftable(*item) {
			npcItem, found := catalog.NpcsItems.Get(code)
			if !found || npcItem.BuyPrice == nil || *npcItem.BuyPrice <= 0 {
				return fmt.Errorf("item is not craftable or available from an NPC: %s", code)
			}

			bankInventory, err := fetchAllBankItems()
			if err != nil {
				return err
			}
			available := 0
			if npcItem.Currency == "gold" {
				bank, err := api.MyBank()
				if err != nil {
					return err
				}
				available = bank.Gold
				for _, character := range characters {
					available += character.Gold
				}
			} else {
				available = bankInventory[npcItem.Currency]
				for _, character := range characters {
					if character.Inventory == nil {
						continue
					}
					for _, invItem := range *character.Inventory {
						if invItem.Code == npcItem.Currency {
							available += invItem.Quantity
						}
					}
				}
			}

			cost := *npcItem.BuyPrice
			quantity := flags.Quantity
			if quantity == 0 {
				quantity = available / cost
			}
			return console.Auto(map[string]interface{}{
				"item":               item.Code,
				"npc":                npcItem.Npc,
				"currency":           npcItem.Currency,
				"cost_per_unit":      cost,
				"currency_available": available,
				"max_exchange":       available / cost,
				"quantity":           quantity,
				"total_cost":         quantity * cost,
			})
		}

		skill := string(*item.Craft.Skill)
		character, compatible := costCraftingCharacter(characters, skill, *item.Craft.Level)
		skillLevel, _ := utils.GetCharacterCraftingSkillLevel(character, skill)
		if character.Name == "" {
			skillLevel = *item.Craft.Level
		}

		bankInventory, err := fetchAllBankItems()
		if err != nil {
			return err
		}

		totalInventory := make(map[string]int)
		if character.Inventory != nil {
			for _, invItem := range *character.Inventory {
				if invItem.Code != "" {
					totalInventory[invItem.Code] = totalInventory[invItem.Code] + invItem.Quantity
				}
			}
		}
		for bCode, amount := range bankInventory {
			totalInventory[bCode] = totalInventory[bCode] + amount
		}

		materialCapacity := craftingCapacity(*item, totalInventory)
		maxCraft := materialCapacity
		quantity, actions := 0, 0
		if compatible {
			quantity, actions, err = craftingTarget(flags.Quantity, maxCraft, *item.Craft.Quantity)
			if err != nil {
				return err
			}
		} else {
			maxCraft = 0
		}
		var bottleneckIngredient string

		ingredientsMap := make(map[string]map[string]int)
		for _, req := range *item.Craft.Items {
			available := totalInventory[req.Code]
			possibleCraft := available / req.Quantity * *item.Craft.Quantity

			ingredientsMap[req.Code] = map[string]int{
				"required":  req.Quantity,
				"available": available,
				"can_craft": possibleCraft,
			}

			if bottleneckIngredient == "" && possibleCraft == materialCapacity {
				bottleneckIngredient = req.Code
			}
		}

		xpPerCraft := CalculateArtifactsXP(*item.Craft.Level, skillLevel, skill, character.Wisdom)
		totalXpGained := actions * xpPerCraft

		output := map[string]interface{}{
			"bottleneck":    bottleneckIngredient,
			"character":     character.Name,
			"ingredients":   ingredientsMap,
			"item":          item.Code,
			"max_craft":     maxCraft,
			"quantity":      quantity,
			"skill":         string(*item.Craft.Skill),
			"xp_per_action": xpPerCraft,
			"xp_total":      totalXpGained,
		}
		if !compatible {
			output["reason"] = costIneligibleReason(character, skill, *item.Craft.Level)
		}

		return console.Auto(output)
	},
}

func costCraftingCharacter(characters []schemas.CharacterSchema, skill string, requiredLevel int) (schemas.CharacterSchema, bool) {
	var lowestCompatible, highestIncompatible schemas.CharacterSchema
	compatibleFound, incompatibleFound := false, false
	for _, character := range characters {
		level, xp := characterSkillProgress(character, skill)
		if level >= requiredLevel {
			selectedLevel, selectedXP := characterSkillProgress(lowestCompatible, skill)
			lowerProgress := level < selectedLevel || level == selectedLevel && xp < selectedXP
			sameProgress := level == selectedLevel && xp == selectedXP
			if !compatibleFound || lowerProgress || sameProgress && character.Name < lowestCompatible.Name {
				lowestCompatible = character
				compatibleFound = true
			}
			continue
		}
		selectedLevel, selectedXP := characterSkillProgress(highestIncompatible, skill)
		higherProgress := level > selectedLevel || level == selectedLevel && xp > selectedXP
		sameProgress := level == selectedLevel && xp == selectedXP
		if !incompatibleFound || higherProgress || sameProgress && character.Name < highestIncompatible.Name {
			highestIncompatible = character
			incompatibleFound = true
		}
	}
	if compatibleFound {
		return lowestCompatible, true
	}
	return highestIncompatible, false
}

func costIneligibleReason(character schemas.CharacterSchema, skill string, requiredLevel int) string {
	if character.Name == "" {
		return fmt.Sprintf("no character can craft %s at level %d", skill, requiredLevel)
	}
	level, _ := utils.GetCharacterCraftingSkillLevel(character, skill)
	return fmt.Sprintf("character %s has %s level %d; requires %d", character.Name, skill, level, requiredLevel)
}

func CalculateArtifactsXP(itemLevel int, playerLevel int, skill string, wisdom int) int {
	var baseXP float64
	var coefficient float64

	if itemLevel < 5 {
		baseXP = 50
		coefficient = 25
	} else if itemLevel >= 5 && itemLevel <= 9 {
		baseXP = 100
		coefficient = 30
	} else if itemLevel >= 10 && itemLevel <= 14 {
		baseXP = 200
		coefficient = 35
	} else if itemLevel >= 15 && itemLevel <= 19 {
		baseXP = 325
		coefficient = 40
	} else if itemLevel >= 20 && itemLevel <= 24 {
		baseXP = 450
		coefficient = 45
	} else if itemLevel >= 25 && itemLevel <= 29 {
		baseXP = 550
		coefficient = 50
	} else if itemLevel >= 30 && itemLevel <= 34 {
		baseXP = 650
		coefficient = 55
	} else if itemLevel >= 35 && itemLevel <= 39 {
		baseXP = 750
		coefficient = 60
	} else if itemLevel >= 40 && itemLevel <= 44 {
		baseXP = 850
		coefficient = 65
	} else {
		baseXP = 1000
		coefficient = 70
	}

	var skillMultiplier float64
	skillMultiplier = 1.0

	switch skill {
	case "fishing", "mining", "woodcutting":
		skillMultiplier = 0.1
	case "cooking":
		skillMultiplier = 0.5
	}

	var levelPenalty float64
	levelPenalty = 1.0

	if playerLevel-itemLevel > 10 {
		levelPenalty = 0.0
	}

	wisdomBonus := 1.0 + (float64(wisdom) * 0.001)

	calculatedXP := (baseXP + (float64(itemLevel) / float64(playerLevel) * coefficient)) * skillMultiplier * levelPenalty * wisdomBonus
	finalXP := math.Round(calculatedXP)

	return int(finalXP)
}

func init() {
	rootCmd.AddCommand(costCmd)
	err := utils.RegisterFlags[costFlags](costCmd)
	if err != nil {
		panic(err)
	}
}
