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

var fillData struct {
	character     schemas.CharacterSchema
	order         schemas.GEOrderSchema
	bankItem      schemas.SimpleItemSchema
	inventoryItem schemas.SimpleItemSchema
}

var fillFlags struct {
	quantity int
}

var fillCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "fill <name> <id>",
	Short: "Fill an existing buy order",
	Long: `Fill an existing buy order.

Arguments:
  name   Name of your character
  id     The id of the buy order you want to fill`,
	ValidArgsFunction: completion.CharacterName(1).Build(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		id := args[1]
		validCharacters := utils.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v",
				name, validCharacters)
		}
		if id == "" {
			return fmt.Errorf("id must not be empty")
		}
		if fillFlags.quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
		orders, err := api.GrandexchangeOrders(id)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			return fmt.Errorf("order %q not found", id)
		}
		fillData.order = orders[0]
		if fillData.order.Type != "buy" {
			return fmt.Errorf("order %q is not a buy order: type %q",
				id, fillData.order.Type)
		}
		if fillFlags.quantity > fillData.order.Quantity {
			return fmt.Errorf(
				"not enough quantity in order %q: required %d, available %d",
				id, fillFlags.quantity, fillData.order.Quantity)
		}
		code := fillData.order.Code
		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.Code == code {
				fillData.bankItem = item
				break
			}
		}
		fillData.character, err = api.Characters(name)
		if err != nil {
			return err
		}
		fillData.inventoryItem = schemas.SimpleItemSchema{}
		if fillData.character.Inventory != nil {
			for _, item := range *fillData.character.Inventory {
				if item.Code == code {
					fillData.inventoryItem = schemas.SimpleItemSchema{
						Code: item.Code, Quantity: item.Quantity}
					break
				}
			}
		}
		totalAvailable := fillData.bankItem.Quantity +
			fillData.inventoryItem.Quantity
		if fillFlags.quantity > totalAvailable {
			return fmt.Errorf("not enough items: required %d, available %d",
				fillFlags.quantity, totalAvailable)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		name := args[0]
		id := args[1]
		code := fillData.order.Code
		totalFilled := 0
		for totalFilled < fillFlags.quantity {
			remaining := fillFlags.quantity - totalFilled
			if fillData.inventoryItem.Quantity == 0 {
				var err error
				fillData.character, err =
					routine.Move(fillData.character, "bank")
				if err != nil {
					return err
				}
				items := routine.GetInventoryItems(fillData.character, nil)
				if len(items) > 0 {
					depositData, err := api.MyActionBankDepositItem(
						fillData.character.Name, items)
					if err != nil {
						return err
					}
					fillData.character = depositData.Character
				}
				withdrawQuantity := min(remaining,
					fillData.character.InventoryMaxItems)
				withdrawItem := schemas.SimpleItemSchema{
					Code: code, Quantity: withdrawQuantity}
				withdrawData, err := api.MyActionBankWithdrawItem(name,
					[]schemas.SimpleItemSchema{withdrawItem})
				if err != nil {
					return err
				}
				fillData.character = withdrawData.Character
				fillData.inventoryItem = schemas.SimpleItemSchema{
					Code: code, Quantity: withdrawQuantity}
			}
			var err error
			fillData.character, err =
				routine.Move(fillData.character, "grand_exchange")
			if err != nil {
				return err
			}
			fill := schemas.GEFillBuyOrderSchema{
				Id:       id,
				Quantity: min(remaining, fillData.inventoryItem.Quantity, 100),
			}
			fillOrderData, err := api.MyActionGrandexchangeFill(name, fill)
			if err != nil {
				return err
			}
			fillData.character = fillOrderData.Character
			fillData.inventoryItem.Quantity -= fillOrderData.Order.Quantity
			totalFilled += fillOrderData.Order.Quantity
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fillCmd)
	fillCmd.Flags().IntVarP(&fillFlags.quantity, "quantity", "q", 0,
		"Item quantity")
}
