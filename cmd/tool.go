package cmd

import (
	_ "embed"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var validSkills = []string{
	"mining",
	"woodcutting",
	"fishing",
	"alchemy",
	"cooking",
	"gearcrafting",
	"jewelrycrafting",
	"weaponcrafting",
}

var toolCmd = &cobra.Command{
	Use:   "tool <name> <code>",
	Short: "Get Best Tools For Skill",
	Long: `Get Best Tools For Skill

Arguments:
  <name>   The name of the character.
  <code>   The skill code.`,
	Args: cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			characters := config.GetCharacters()
			return characters, cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			return validSkills, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		if !slices.Contains(validSkills, code) {
			return fmt.Errorf("invalid skill code: %s (valid options: %s)",
				code, strings.Join(validSkills, ", "))
		}

		character, err := api.Characters(name)
		if err != nil {
			return fmt.Errorf("failed to get character: %w", err)
		}

		console.Auto(
			GetBestEquipmentForSkill(character, code, database.GetItems()))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(toolCmd)
}

func getCharacterSkillLevel(character schemas.CharacterSchema, skill string) int {
	switch skill {
	case "mining":
		return character.MiningLevel
	case "woodcutting":
		return character.WoodcuttingLevel
	case "fishing":
		return character.FishingLevel
	case "alchemy":
		return character.AlchemyLevel
	case "cooking":
		return character.CookingLevel
	case "gearcrafting":
		return character.GearcraftingLevel
	case "jewelrycrafting":
		return character.JewelrycraftingLevel
	case "weaponcrafting":
		return character.WeaponcraftingLevel
	default:
		return character.Level
	}
}

func GetBestEquipmentForSkill(character schemas.CharacterSchema, skill string, items []schemas.ItemSchema) []string {
	var filteredItems []schemas.ItemSchema
	characterLevel := getCharacterSkillLevel(character, skill)

	for _, item := range items {
		if item.Type == "weapon" {
			if item.Subtype == "tool" {
				var hasSkillEffect bool

				if item.Effects != nil {
					for _, effect := range *item.Effects {
						if effect.Code == skill {
							hasSkillEffect = true
						}
					}
				}

				if hasSkillEffect {
					var canEquip bool
					canEquip = true

					if item.Conditions != nil {
						for _, condition := range *item.Conditions {
							if condition.Code == "level" {
								if condition.Operator == "gt" {
									if characterLevel <= condition.Value {
										canEquip = false
									}
								}
							}

							if condition.Code == skill+"_level" {
								if condition.Operator == "gt" {
									if characterLevel <= condition.Value {
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
	}

	sort.Slice(filteredItems, func(i, j int) bool {
		var valueI int
		var valueJ int

		if filteredItems[i].Effects != nil {
			for _, effect := range *filteredItems[i].Effects {
				if effect.Code == skill {
					valueI = effect.Value
				}
			}
		}

		if filteredItems[j].Effects != nil {
			for _, effect := range *filteredItems[j].Effects {
				if effect.Code == skill {
					valueJ = effect.Value
				}
			}
		}

		return valueI < valueJ
	})

	var result []string

	for _, item := range filteredItems {
		result = append(result, item.Code)
	}

	return result
}
