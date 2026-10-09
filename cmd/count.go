package cmd

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type countFlags struct {
	Type       []string `flag:"type" shorthand:"t" desc:"Item types to include"`
	Subtype    []string `flag:"subtype" desc:"Item subtypes to include"`
	MinLevel   int      `flag:"min-level" default:"1" desc:"Minimum item level"`
	MaxLevel   int      `flag:"max-level" default:"50" desc:"Maximum item level"`
	Level      bool     `flag:"level" desc:"Include item levels in the output"`
	Skill      []string `flag:"skill" shorthand:"s" desc:"Crafting skills to include"`
	Tradeable  bool     `flag:"tradeable" desc:"Only include tradeable items"`
	Recyclable bool     `flag:"recyclable" desc:"Only include recyclable items"`
}

var countCmd = &cobra.Command{
	Args:  cobra.ArbitraryArgs,
	Use:   "count [code...]",
	Short: "Show the total quantity of items in the account",
	Long: `Show the total quantity of items in the account

Arguments:
  code   The code of the item (or 'gold' for gold).`,
	ValidArgsFunction: completion.Custom(0, func() []string {
		return append(catalog.Items().Keys(), "gold")
	}).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[countFlags](cmd)
		if err != nil {
			return err
		}
		err = countValidate(args, flags)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return countRun(args, flags)
	},
}

func countValidate(args []string, flags countFlags) error {
	for _, code := range args {
		if code == "gold" {
			continue
		}
		_, found := catalog.Items().Get(code)
		if !found {
			return fmt.Errorf("item %s not found", code)
		}
	}
	for _, itemType := range flags.Type {
		if !slices.Contains(catalog.Enums()["ItemType"], itemType) {
			return fmt.Errorf("invalid item type %q", itemType)
		}
	}
	for _, subtype := range flags.Subtype {
		if !slices.Contains(countSubtypes(), subtype) {
			return fmt.Errorf("invalid item subtype %q", subtype)
		}
	}
	for _, skill := range flags.Skill {
		if !slices.Contains(catalog.Enums()["CraftSkill"], skill) {
			return fmt.Errorf("invalid crafting skill %q", skill)
		}
	}
	err := validateGameLevel(flags.MinLevel)
	if err != nil {
		return err
	}
	err = validateGameLevel(flags.MaxLevel)
	if err != nil {
		return err
	}
	if flags.MinLevel > flags.MaxLevel {
		return fmt.Errorf("minimum level cannot exceed maximum level")
	}
	return nil
}

func countRun(args []string, flags countFlags) error {
	bank, err := api.MyBank()
	if err != nil {
		return err
	}
	items, err := api.MyBankItems()
	if err != nil {
		return err
	}
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return err
	}

	output := map[string][]map[string]int{}
	slotQuantities := make([]map[string]int, len(characters))
	for index, character := range characters {
		slotQuantities[index] = countInSlots(character)
	}
	codes := countCodes(args, items, characters, slotQuantities)

	for _, code := range codes {
		level := 0
		if code != "gold" {
			item, found := catalog.Items().Get(code)
			if !found || !countItemMatches(*item, flags) {
				continue
			}
			level = item.Level
		}
		result := []map[string]int{}

		if code == "gold" {
			if bank.Gold > 0 {
				result = append(result, map[string]int{"bank": bank.Gold})
			}
		} else {
			for _, item := range items {
				if item.Code == code {
					result = append(result, map[string]int{
						"bank": item.Quantity,
					})
				}
			}
		}

		for index, character := range characters {
			if code == "gold" {
				if character.Gold > 0 {
					result = append(result, map[string]int{
						character.Name: character.Gold,
					})
					continue
				}
			}
			inventory := []schemas.InventorySlotSchema{}
			if character.Inventory != nil {
				inventory = *character.Inventory
			}
			qty := 0
			for _, item := range inventory {
				if item.Code == code {
					qty += item.Quantity
				}
			}
			qty += slotQuantities[index][code]
			if qty > 0 {
				result = append(result, map[string]int{character.Name: qty})
			}
		}

		total := 0
		for _, entry := range result {
			for _, qty := range entry {
				total += qty
			}
		}
		if flags.Level && code != "gold" {
			result = append(result, map[string]int{"level": level})
		}
		result = append(result, map[string]int{"total": total})

		output[code] = result
	}

	return console.Auto(output)
}

func countCodes(args []string, bankItems []schemas.SimpleItemSchema, characters []schemas.CharacterSchema, slotQuantities []map[string]int) []string {
	if len(args) > 0 {
		return args
	}
	codes := map[string]bool{}
	for _, item := range bankItems {
		if item.Code != "" && item.Quantity > 0 {
			codes[item.Code] = true
		}
	}
	for index, character := range characters {
		if character.Inventory != nil {
			for _, item := range *character.Inventory {
				if item.Code != "" && item.Quantity > 0 {
					codes[item.Code] = true
				}
			}
		}
		for code, quantity := range slotQuantities[index] {
			if quantity > 0 {
				codes[code] = true
			}
		}
	}
	result := make([]string, 0, len(codes))
	for code := range codes {
		result = append(result, code)
	}
	sort.Strings(result)
	return result
}

func countItemMatches(item schemas.ItemSchema, flags countFlags) bool {
	if len(flags.Type) > 0 && !slices.Contains(flags.Type, item.Type) {
		return false
	}
	if len(flags.Subtype) > 0 && !slices.Contains(flags.Subtype, item.Subtype) {
		return false
	}
	if item.Level < flags.MinLevel || item.Level > flags.MaxLevel {
		return false
	}
	if flags.Tradeable && !item.Tradeable {
		return false
	}
	if flags.Recyclable && (item.Recyclable == nil || !*item.Recyclable) {
		return false
	}
	if len(flags.Skill) == 0 {
		return true
	}
	return item.Craft != nil && item.Craft.Skill != nil && slices.Contains(flags.Skill, *item.Craft.Skill)
}

func countSubtypes() []string {
	subtypes := map[string]bool{}
	for _, item := range catalog.Items().All() {
		if item.Subtype != "" {
			subtypes[item.Subtype] = true
		}
	}
	result := make([]string, 0, len(subtypes))
	for subtype := range subtypes {
		result = append(result, subtype)
	}
	sort.Strings(result)
	return result
}

func countInSlots(character schemas.CharacterSchema) map[string]int {
	quantities := map[string]int{}
	v := reflect.ValueOf(character)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		if !strings.HasSuffix(fieldName, "Slot") || t.Field(i).Type.Kind() != reflect.String {
			continue
		}
		code := v.Field(i).String()
		if code == "" {
			continue
		}
		quantityField := v.FieldByName(fieldName + "Quantity")
		if quantityField.IsValid() && quantityField.Kind() == reflect.Int {
			quantities[code] += int(quantityField.Int())
		} else {
			quantities[code]++
		}
	}
	return quantities
}

func init() {
	rootCmd.AddCommand(countCmd)
	err := utils.RegisterFlags[countFlags](countCmd)
	if err != nil {
		panic(err)
	}
	err = countCmd.RegisterFlagCompletionFunc("type", func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return catalog.Enums()["ItemType"], cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = countCmd.RegisterFlagCompletionFunc("subtype", func(
		cmd *cobra.Command,
		args []string,
		toComplete string,
	) ([]string, cobra.ShellCompDirective) {
		return countSubtypes(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = countCmd.RegisterFlagCompletionFunc("skill", completion.StringSlice(func() []string {
		return catalog.Enums()["CraftSkill"]
	}))
	if err != nil {
		panic(err)
	}
}
