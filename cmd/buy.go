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

var buyData struct {
	character schemas.CharacterSchema
	order     schemas.GEOrderSchema
	bankGold  int
}

var buyFlags struct {
	quantity int
}

var buyCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "buy <name> <id>",
	Short: "Buy an existing sell order",
	Long: `Buy an existing sell order

Arguments:
  name   Name of your character.
  id     The id of the sell order you want to buy.`,
	ValidArgsFunction: completion.CharacterName(1).Build(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		id := args[1]
		validCharacters := utils.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
		}
		if id == "" {
			return fmt.Errorf("id must not be empty")
		}
		if buyFlags.quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
		order, err := api.GrandexchangeOrder(id)
		if err != nil {
			return err
		}
		buyData.order = order
		if buyData.order.Type != "sell" {
			return fmt.Errorf("order %q is not a sell order: type %q", id, buyData.order.Type)
		}
		if buyFlags.quantity > buyData.order.Quantity {
			return fmt.Errorf("not enough quantity in order %q: required %d, available %d", id, buyFlags.quantity, buyData.order.Quantity)
		}
		bank, err := api.MyBank()
		if err != nil {
			return err
		}
		buyData.bankGold = bank.Gold
		buyData.character, err = api.Characters(name)
		if err != nil {
			return err
		}
		totalGold := buyData.character.Gold + buyData.bankGold
		requiredGold := buyFlags.quantity * buyData.order.Price
		if totalGold < requiredGold {
			return fmt.Errorf("not enough gold: required %d, available %d", requiredGold, totalGold)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		name := args[0]
		id := args[1]
		totalBought := 0
		for totalBought < buyFlags.quantity {
			remaining := buyFlags.quantity - totalBought
			quantity := min(remaining, buyData.character.InventoryMaxItems, 100)
			cost := quantity * buyData.order.Price

			needBank := false
			if buyData.character.Inventory != nil {
				for _, item := range *buyData.character.Inventory {
					if item.Code != "" {
						needBank = true
						break
					}
				}
			}
			if needBank || buyData.character.Gold < cost {
				var err error
				buyData.character, err = routine.Move(buyData.character, "bank")
				if err != nil {
					return err
				}
				items := routine.GetInventoryItems(buyData.character, nil)
				if len(items) > 0 {
					depositData, err := api.MyActionBankDepositItem(name, items)
					if err != nil {
						return err
					}
					buyData.character = depositData.Character
				}
				if buyData.character.Gold < cost {
					goldData, err := api.MyActionBankWithdrawGold(name, cost-buyData.character.Gold)
					if err != nil {
						return err
					}
					buyData.character = goldData.Character
				}
			}

			var err error
			buyData.character, err = routine.Move(buyData.character, "grand_exchange")
			if err != nil {
				return err
			}

			buy := schemas.GEBuyOrderSchema{Id: id, Quantity: quantity}
			buyOrderData, err := api.MyActionGrandexchangeBuy(name, buy)
			if err != nil {
				return err
			}
			buyData.character = buyOrderData.Character
			totalBought += buyOrderData.Order.Quantity
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buyCmd)
	buyCmd.Flags().IntVarP(&buyFlags.quantity, "quantity", "q", 0, "Item quantity")
}
