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

var elements = []string{"fire", "water", "earth", "air"}

var weaponCmd = &cobra.Command{
	Use:   "weapon <name> <monster>",
	Short: "Get Best Weapon For Monster",
	Long: `Get Best Weapon For Monster

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
			GetBestWeaponForMonster(character, monster, database.GetItems()))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(weaponCmd)
}

func GetBestWeaponForMonster(character schemas.CharacterSchema, monster schemas.MonsterSchema, items []schemas.ItemSchema) []string {
	var filteredItems []schemas.ItemSchema

	for _, item := range items {
		if item.Type == "weapon" {
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
		return calculateWeaponDamage(filteredItems[i], monster) >
			calculateWeaponDamage(filteredItems[j], monster)
	})

	var result []string
	for _, item := range filteredItems {
		result = append(result, item.Code)
	}

	return result
}

func calculateWeaponDamage(item schemas.ItemSchema, monster schemas.MonsterSchema) float64 {
	var totalDamage float64
	totalDamage = 0.0

	if item.Effects == nil {
		return totalDamage
	}

	for _, effect := range *item.Effects {
		for _, element := range elements {
			if effect.Code == "attack_"+element {
				var resistance int
				resistance = getMonsterResistance(monster, element)

				var damage float64
				damage = float64(effect.Value) *
					(1.0 - (float64(resistance) / 100.0))
				totalDamage = totalDamage + damage
			}
		}
	}

	return totalDamage
}

func getMonsterResistance(monster schemas.MonsterSchema, element string) int {
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
