package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/surplus"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type missingFlags struct {
	Name          []string `flag:"name" shorthand:"n" desc:"Character names to analyze"`
	EquipmentType []string `flag:"type" shorthand:"t" desc:"Equipment types to filter"`
	Craftable     bool     `flag:"craftable" desc:"Only include items craftable by a character"`
}

var missingOptions missingFlags

var missingCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "missing",
	Short: "Determine missing progression gear for characters",

	ValidArgsFunction: completion.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		missingOptions, err = utils.ReadFlags[missingFlags](cmd)
		if err != nil {
			return err
		}
		validCharacters := cache.GetCharacters()
		for _, name := range missingOptions.Name {
			if !slices.Contains(validCharacters, name) {
				return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
			}
		}
		for _, equipmentType := range missingOptions.EquipmentType {
			if !slices.Contains(database.EquipmentTypes, equipmentType) {
				return fmt.Errorf("invalid equipment type specified: %s", equipmentType)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return executeMissing(missingOptions)
	},
}

func executeMissing(flags missingFlags) error {
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return err
	}
	var charactersToProcess []schemas.CharacterSchema
	if len(flags.Name) > 0 {
		for _, char := range characters {
			hasChar := slices.Contains(flags.Name, char.Name)
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
	ownedItems := make(map[string]int)
	for _, item := range bankItems {
		ownedItems[item.Code] += item.Quantity
	}
	for _, character := range charactersToProcess {
		if character.Inventory != nil {
			for _, item := range *character.Inventory {
				ownedItems[item.Code] += item.Quantity
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
			character.RuneSlot,
			character.ShieldSlot,
			character.Utility1Slot,
			character.Utility2Slot,
			character.WeaponSlot,
		}
		for _, item := range slots {
			if item != "" {
				ownedItems[item]++
			}
		}
	}
	var targets []string
	if len(flags.EquipmentType) > 0 {
		targets = flags.EquipmentType
	} else {
		targets = database.EquipmentTypes
	}
	allItems := database.Items().All()
	requiredItems := make(map[string]int)
	craftableItems := make(map[string]bool)
	if flags.Craftable {
		for _, item := range allItems {
			for _, character := range characters {
				if canCraft(character, *item) {
					craftableItems[item.Code] = true
					break
				}
			}
		}
	}
	for _, character := range charactersToProcess {
		var candidates []schemas.ItemSchema
		for _, item := range allItems {
			hasStaff := ownedItems["wooden_staff"] > 0
			if item.Code == "wooden_stick" && hasStaff {
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
				if !hasRequiredSkillLevel(character, *item) {
					continue
				}
			} else if item.Level > character.Level {
				continue
			}
			candidates = append(candidates, *item)
		}
		for _, item := range surplus.NonDominated(candidates, character) {
			requiredItems[item.Code] += requiredEquipmentQuantity(item)
		}
	}
	missingItems := make(map[string]int)
	for code, required := range requiredItems {
		if flags.Craftable && !craftableItems[code] {
			continue
		}
		missing := required - ownedItems[code]
		if missing > 0 {
			missingItems[code] = missing
		}
	}
	return console.Auto(missingItems)
}

func requiredEquipmentQuantity(item schemas.ItemSchema) int {
	if item.Type == "ring" {
		return 2
	}
	return 1
}

func canCraft(character schemas.CharacterSchema, item schemas.ItemSchema) bool {
	if item.Craft == nil || item.Craft.Skill == nil || item.Craft.Level == nil {
		return false
	}
	level, exists := utils.GetCharacterCraftingSkillLevel(character, string(*item.Craft.Skill))
	return exists && level >= *item.Craft.Level
}

func hasRequiredSkillLevel(character schemas.CharacterSchema, item schemas.ItemSchema) bool {
	if item.Conditions == nil {
		return true
	}
	conditions := *item.Conditions
	for _, condition := range conditions {
		currentLevel, _ := utils.GetCharacterConditionLevel(character, condition.Code)
		if currentLevel == 0 {
			continue
		}
		satisfied := true
		switch condition.Operator {
		case schemas.Gt:
			satisfied = currentLevel > condition.Value
		case schemas.Eq:
			satisfied = currentLevel == condition.Value
		case schemas.Lt:
			satisfied = currentLevel < condition.Value
		case schemas.Ne:
			satisfied = currentLevel != condition.Value
		}
		if !satisfied {
			return false
		}
	}
	return true
}

func init() {
	rootCmd.AddCommand(missingCmd)
	err := utils.RegisterFlags[missingFlags](missingCmd)
	if err != nil {
		panic(err)
	}
	err = missingCmd.RegisterFlagCompletionFunc("name", completion.StringSlice(cache.GetCharacters))
	if err != nil {
		panic(err)
	}
	err = missingCmd.RegisterFlagCompletionFunc("type", completion.StringSlice(func() []string {
		return database.EquipmentTypes
	}))
	if err != nil {
		panic(err)
	}
}
