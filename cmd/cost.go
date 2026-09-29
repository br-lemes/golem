package cmd

import (
	"fmt"
	"math"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type costFlags struct {
	Quantity int `flag:"quantity" shorthand:"q" desc:"Amount of items to craft (0 for maximum available)"`
}

var costCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "cost <name> <code>",
	Short: "Calculate resources and XP for crafting an item",
	Long: `Calculate resources and XP for crafting an item

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

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
		if !isCraftable(*item) && flags.Quantity != 0 {
			return fmt.Errorf("quantity is only supported for craftable items")
		}
		cmd.SilenceUsage = true

		character, err := api.Characters(name)
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
				available = bank.Gold + character.Gold
			} else {
				available = bankInventory[npcItem.Currency]
				for _, invItem := range *character.Inventory {
					if invItem.Code == npcItem.Currency {
						available += invItem.Quantity
					}
				}
			}

			cost := *npcItem.BuyPrice
			return console.Auto(map[string]interface{}{
				"item":               item.Code,
				"npc":                npcItem.Npc,
				"currency":           npcItem.Currency,
				"cost_per_unit":      cost,
				"currency_available": available,
				"max_exchange":       available / cost,
			})
		}

		skillLevel, _ := utils.GetCharacterCraftingSkillLevel(character, string(*item.Craft.Skill))
		if skillLevel < *item.Craft.Level {
			return fmt.Errorf("character %s level too low. Required: %d, Current: %d", name, *item.Craft.Level, skillLevel)
		}

		bankInventory, err := fetchAllBankItems()
		if err != nil {
			return err
		}

		totalInventory := make(map[string]int)
		for _, invItem := range *character.Inventory {
			if invItem.Code != "" {
				totalInventory[invItem.Code] = totalInventory[invItem.Code] + invItem.Quantity
			}
		}
		for bCode, amount := range bankInventory {
			totalInventory[bCode] = totalInventory[bCode] + amount
		}

		maxCraft := craftingCapacity(*item, totalInventory)
		quantity, actions, err := craftingTarget(flags.Quantity, maxCraft, *item.Craft.Quantity)
		if err != nil {
			return err
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

			if bottleneckIngredient == "" && possibleCraft == maxCraft {
				bottleneckIngredient = req.Code
			}
		}

		xpPerCraft := CalculateArtifactsXP(*item.Craft.Level, skillLevel, string(*item.Craft.Skill), character.Wisdom)
		totalXpGained := actions * xpPerCraft

		output := map[string]interface{}{
			"bottleneck":    bottleneckIngredient,
			"ingredients":   ingredientsMap,
			"item":          item.Code,
			"max_craft":     maxCraft,
			"quantity":      quantity,
			"skill":         string(*item.Craft.Skill),
			"xp_per_action": xpPerCraft,
			"xp_total":      totalXpGained,
		}

		return console.Auto(output)
	},
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
