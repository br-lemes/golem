package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var sellData struct {
	character     schemas.CharacterSchema
	bankItem      schemas.SimpleItemSchema
	inventoryItem schemas.SimpleItemSchema
}

var sellFlags struct {
	price    int
	quantity int
}

var sellCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "sell <name> <code>",
	Short: "Create a new sell order",
	Long: `Create a new sell order

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).BankItem(1).Build(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]
		validCharacters := utils.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
		}
		validItems := completion.GetBankItems()
		if !slices.Contains(validItems, code) {
			return fmt.Errorf("invalid item %q: allowed values are %v", code, validItems)
		}
		if sellFlags.price <= 0 {
			return fmt.Errorf("price must be greater than 0")
		}
		if sellFlags.quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.Code == code {
				sellData.bankItem = item
				break
			}
		}
		sellData.character, err = api.Characters(name)
		if err != nil {
			return err
		}
		sellData.inventoryItem = schemas.SimpleItemSchema{}
		if sellData.character.Inventory != nil {
			for _, item := range *sellData.character.Inventory {
				if item.Code == code {
					sellData.inventoryItem = schemas.SimpleItemSchema{
						Code:     item.Code,
						Quantity: item.Quantity,
					}
					break
				}
			}
		}
		totalAvailable := sellData.bankItem.Quantity + sellData.inventoryItem.Quantity
		if sellFlags.quantity > totalAvailable {
			return fmt.Errorf("not enough items: required %d, available %d", sellFlags.quantity, totalAvailable)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		name := args[0]
		code := args[1]
		totalSold := 0
		for totalSold < sellFlags.quantity {
			remaining := sellFlags.quantity - totalSold
			if sellData.inventoryItem.Quantity == 0 {
				var err error
				sellData.character, err = routine.Move(sellData.character, "bank")
				if err != nil {
					return err
				}
				items := routine.GetInventoryItems(sellData.character, nil)
				if len(items) > 0 {
					depositData, err := api.MyActionBankDepositItem(sellData.character.Name, items)
					if err != nil {
						return err
					}
					sellData.character = depositData.Character
				}
				withdrawQuantity := min(remaining, sellData.character.InventoryMaxItems)
				withdrawData, err := api.MyActionBankWithdrawItem(name, []schemas.SimpleItemSchema{
					{Code: code, Quantity: withdrawQuantity},
				})
				if err != nil {
					return err
				}
				sellData.character = withdrawData.Character
				sellData.inventoryItem = schemas.SimpleItemSchema{
					Code:     code,
					Quantity: withdrawQuantity,
				}
			}
			var err error
			sellData.character, err = routine.Move(sellData.character, "grand_exchange")
			if err != nil {
				return err
			}
			order := schemas.GEOrderCreationSchema{
				Code:     code,
				Price:    sellFlags.price,
				Quantity: min(remaining, sellData.inventoryItem.Quantity, 100),
			}
			orderData, err := api.MyActionGrandexchangeCreateSellOrder(name, order)
			if err != nil {
				return err
			}
			sellData.character = orderData.Character
			sellData.inventoryItem.Quantity -= orderData.Order.Quantity
			totalSold += orderData.Order.Quantity
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sellCmd)
	sellCmd.Flags().IntVarP(&sellFlags.price, "price", "p", 0, "Item price per unit")
	sellCmd.Flags().IntVarP(&sellFlags.quantity, "quantity", "q", 0, "Item quantity")
}
