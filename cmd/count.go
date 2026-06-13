package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
	Use:   "count <code>",
	Short: "Show the total quantity of an item in the account",
	Long: `Show the total quantity of an item in the account

Arguments:
  code      The code of the item.`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return database.GetItemCodes(), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		result := []map[string]int{}
		_, found := database.GetItem(code)
		if !found {
			return fmt.Errorf("item %s not found", code)
		}

		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.Code == code {
				result = append(result, map[string]int{"bank": item.Quantity})
			}
		}
		characters, err := api.AccountsCharacters(config.Account)
		if err != nil {
			return err
		}
		for _, character := range characters {
			inventory := []schemas.InventorySlot{}
			if character.Inventory != nil {
				inventory = *character.Inventory
			}
			qty := 0
			for _, item := range inventory {
				if item.Code == code {
					qty += item.Quantity
				}
			}
			qty += countInSlots(character, code)
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
		result = append(result, map[string]int{"total": total})
		return console.Auto(result)
	},
}

func init() {
	rootCmd.AddCommand(countCmd)
}

func countInSlots(character schemas.CharacterSchema, code string) int {
	total := 0
	v := reflect.ValueOf(character)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		if !strings.HasSuffix(fieldName, "Slot") ||
			t.Field(i).Type.Kind() != reflect.String {
			continue
		}
		if v.Field(i).String() != code {
			continue
		}
		quantityField := v.FieldByName(fieldName + "Quantity")
		if quantityField.IsValid() && quantityField.Kind() == reflect.Int {
			total += int(quantityField.Int())
		} else {
			total += 1
		}
	}
	return total
}
