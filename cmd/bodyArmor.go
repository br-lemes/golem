package cmd

import (
	"fmt"
	"sort"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var bodyArmorCmd = &cobra.Command{
	Use:   "bodyArmor <name> <monster>",
	Short: "Get Best Body Armor For Monster",
	Long: `Get Best Body Armor For Monster

Arguments:
  name   Name of your character.
  code   The code of the monster.`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			characters := config.GetCharacters()
			return characters, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			monsters := database.GetMonsterCodes()
			return monsters, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		monsterCode := args[1]

		monster, found := database.GetMonster(monsterCode)
		if !found {
			return fmt.Errorf("monster not found")
		}

		character, err := api.Characters(name)
		if err != nil {
			return fmt.Errorf("failed to get character: %w", err)
		}

		console.Auto(
			GetBestBodyArmorForMonster(character, monster, database.GetItems()))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bodyArmorCmd)
}

func GetBestBodyArmorForMonster(character schemas.CharacterSchema, monster schemas.MonsterSchema, items []schemas.ItemSchema) []string {
	var filteredItems []schemas.ItemSchema

	for _, item := range items {
		if item.Type == "body_armor" {
			if item.Subtype == "" {
				var canEquip bool
				canEquip = true

				if item.Conditions != nil {
					for _, condition := range *item.Conditions {
						if condition.Code == "level" {
							if condition.Operator == "gt" {
								if character.Level <= condition.Value {
									canEquip = false
								}
							}
						}
					}
				}

				if canEquip {
					filteredItems = append(filteredItems, item)
				}
			}
		}
	}

	sort.Slice(filteredItems, func(i, j int) bool {
		return evaluateBodyArmorValue(filteredItems[i], monster) >
			evaluateBodyArmorValue(filteredItems[j], monster)
	})

	var result []string
	for _, item := range filteredItems {
		result = append(result, item.Code)
	}

	return result
}

func evaluateBodyArmorValue(item schemas.ItemSchema, monster schemas.MonsterSchema) float64 {
	var score float64
	score = 0.0

	if item.Effects == nil {
		return score
	}

	for _, effect := range *item.Effects {
		if effect.Code == "hp" {
			score = score + float64(effect.Value)
		}

		for _, element := range elements {
			if effect.Code == "res_"+element {
				var monsterAttack int
				monsterAttack = getMonsterAttackValue(monster, element)
				score = score + (float64(monsterAttack) * (float64(effect.Value) / 100.0))
			}

			if effect.Code == "dmg_"+element {
				var monsterResistance int
				monsterResistance = getMonsterResistanceValue(monster, element)
				score = score + (float64(effect.Value) * (1.0 - (float64(monsterResistance) / 100.0)))
			}
		}

		if effect.Code == "critical_strike" {
			score = score + (float64(effect.Value) * 0.5)
		}

		if effect.Code == "haste" {
			score = score + (float64(effect.Value) * 0.5)
		}
	}

	return score
}

func getMonsterAttackValue(monster schemas.MonsterSchema, element string) int {
	switch element {
	case "fire":
		return monster.AttackFire
	case "water":
		return monster.AttackWater
	case "earth":
		return monster.AttackEarth
	case "air":
		return monster.AttackAir
	default:
		return 0
	}
}

func getMonsterResistanceValue(monster schemas.MonsterSchema, element string) int {
	switch element {
	case "fire":
		return monster.ResFire
	case "water":
		return monster.ResWater
	case "earth":
		return monster.ResEarth
	case "air":
		return monster.ResAir
	default:
		return 0
	}
}
