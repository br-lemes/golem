package cmd

import (
	"errors"
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var missingCmd = &cobra.Command{
	Use:   "missing [<name> <type> | --name name,name --type type,type]",
	Short: "Determine missing progression gear for characters",
	Long: `Determine missing progression gear for characters

Arguments:
  name   The name of the character.
  type   The type of equipment.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		err := cobra.ExactArgs(2)(cmd, args)
		if err != nil {
			return errors.New("positional usage requires exactly 2 arguments: <name> <type>. Otherwise, use --name and --type flags")
		}
		err = validateEquipmentType(args[1])
		if err != nil {
			return err
		}
		return nil
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
		}
		if len(args) == 1 {
			return database.EquipmentTypes, cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		characterNames, _ := cmd.Flags().GetStringSlice("name")
		equipmentTypes, _ := cmd.Flags().GetStringSlice("type")
		if len(args) == 2 {
			characterNames = []string{args[0]}
			equipmentTypes = []string{args[1]}
		}
		for _, t := range equipmentTypes {
			err := validateEquipmentType(t)
			if err != nil {
				return err
			}
		}
		return executeMissing(characterNames, equipmentTypes)
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)
	missingCmd.Flags().StringSliceP("name", "n", []string{}, "Character names to analyze")
	missingCmd.Flags().StringSliceP("type", "t", []string{}, "Equipment types to filter")
	err := missingCmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.GetCharacters(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = missingCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return database.EquipmentTypes, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}

func validateEquipmentType(t string) error {
	if !slices.Contains(database.EquipmentTypes, t) {
		return fmt.Errorf("invalid equipment type specified: %s", t)
	}
	return nil
}

func executeMissing(names []string, equipTypes []string) error {
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return err
	}
	var charactersToProcess []schemas.CharacterSchema
	if len(names) > 0 {
		for _, char := range characters {
			hasChar := slices.Contains(names, char.Name)
			if hasChar {
				charactersToProcess = append(charactersToProcess, char)
			}
		}
	} else {
		charactersToProcess = characters
	}
	bankItems, err := api.MyBankItems()
	if err != nil {
		return err
	}
	ownedItems := make(map[string]bool)
	for _, item := range bankItems {
		ownedItems[item.Code] = true
	}
	var targets []string
	if len(equipTypes) > 0 {
		targets = equipTypes
	} else {
		targets = database.EquipmentTypes
	}
	allItems := database.GetItems()
	var filteredCodes []string
	for _, character := range charactersToProcess {
		if character.Inventory != nil {
			inventorySlots := *character.Inventory
			for _, item := range inventorySlots {
				ownedItems[item.Code] = true
			}
		}
		slots := [...]string{
			character.AmuletSlot,
			character.Artifact1Slot,
			character.Artifact2Slot,
			character.Artifact3Slot,
			character.BagSlot,
			character.BodyArmorSlot,
			character.BootsSlot,
			character.HelmetSlot,
			character.LegArmorSlot,
			character.Ring1Slot,
			character.Ring2Slot,
			character.ShieldSlot,
			character.Utility1Slot,
			character.Utility2Slot,
			character.WeaponSlot,
		}
		for _, slotItem := range slots {
			if slotItem != "" {
				ownedItems[slotItem] = true
			}
		}
		for _, item := range allItems {
			hasStaff := ownedItems["wooden_staff"]
			if item.Code == "wooden_stick" && hasStaff {
				continue
			}
			hasItem := ownedItems[item.Code]
			if hasItem {
				continue
			}
			matchedType := false
			isTool := false
			for _, target := range targets {
				if target == "tool" {
					if item.Type == "weapon" && item.Subtype == "tool" {
						matchedType = true
						isTool = true
					}
				} else {
					if item.Type == target && item.Subtype != "tool" {
						matchedType = true
					}
				}
			}
			if !matchedType {
				continue
			}
			if isTool {
				canEquip := hasRequiredSkillLevel(character, item)
				if !canEquip {
					continue
				}
			} else {
				if item.Level > character.Level {
					continue
				}
			}
			if !slices.Contains(filteredCodes, item.Code) {
				filteredCodes = append(filteredCodes, item.Code)
			}
		}
	}
	return console.Auto(filteredCodes)
}

func hasRequiredSkillLevel(character schemas.CharacterSchema, item schemas.ItemSchema) bool {
	if item.Conditions == nil {
		return true
	}
	conditions := *item.Conditions
	for _, condition := range conditions {
		currentLevel := getSkillLevel(character, condition.Code)
		if currentLevel > 0 && currentLevel < condition.Value {
			return false
		}
	}
	return true
}

func getSkillLevel(character schemas.CharacterSchema, code string) int {
	switch code {
	case "alchemy_level":
		return character.AlchemyLevel
	case "fishing_level":
		return character.FishingLevel
	case "mining_level":
		return character.MiningLevel
	case "woodcutting_level":
		return character.WoodcuttingLevel
	default:
		return 0
	}
}
