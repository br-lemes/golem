package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var relistData struct {
	character schemas.CharacterSchema
	orders    []schemas.GEOrderSchema
}

var relistFlags struct {
	price    int
	quantity int
}

var relistCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "relist <name> <code>",
	Short: "Relist existing sell orders for an item at a new price",
	Long: `Relist existing sell orders for an item at a new price

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).Tradeable(1).Build(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]
		validCharacters := cache.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
		}
		_, found := database.Items().Tradeables().Get(code)
		if !found {
			return fmt.Errorf("item %q not tradeable or not found", code)
		}
		if relistFlags.price <= 0 {
			return fmt.Errorf("price must be greater than 0")
		}
		if relistFlags.quantity < 0 {
			return fmt.Errorf("quantity must be greater than or equal to 0")
		}
		orders, err := api.MyGrandexchangeOrders(api.MyGrandexchangeOrdersOptions{
			Code: code,
			Type: "sell",
		})
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			return fmt.Errorf("no sell orders found for item %q", code)
		}
		relistData.orders = orders
		relistData.character, err = api.Characters(name)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		name := args[0]
		code := args[1]

		totalToRelist := 0
		if relistFlags.quantity > 0 {
			totalToRelist = relistFlags.quantity
		} else {
			for _, order := range relistData.orders {
				totalToRelist += order.Quantity
			}
		}

		totalProcessed := 0
		for _, order := range relistData.orders {
			if totalProcessed >= totalToRelist {
				break
			}
			remaining := totalToRelist - totalProcessed
			qtyToCancel := min(remaining, order.Quantity)

			var err error
			relistData.character, err = routine.Move(relistData.character, "grand_exchange")
			if err != nil {
				return err
			}

			cancel := schemas.GECancelOrderSchema{Id: order.Id}
			cancelData, err := api.MyActionGrandexchangeCancel(name, cancel)
			if err != nil {
				return err
			}
			relistData.character = cancelData.Character

			oldQtyToRecreate := order.Quantity - qtyToCancel
			if oldQtyToRecreate > 0 {
				oldOrder := schemas.GEOrderCreationSchema{
					Code:     code,
					Price:    order.Price,
					Quantity: oldQtyToRecreate,
				}
				createOldData, err := api.MyActionGrandexchangeCreateSellOrder(name, oldOrder)
				if err != nil {
					return err
				}
				relistData.character = createOldData.Character
			}

			newOrder := schemas.GEOrderCreationSchema{
				Code:     code,
				Price:    relistFlags.price,
				Quantity: qtyToCancel,
			}
			createNewData, err := api.MyActionGrandexchangeCreateSellOrder(name, newOrder)
			if err != nil {
				return err
			}
			relistData.character = createNewData.Character
			totalProcessed += qtyToCancel
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(relistCmd)
	relistCmd.Flags().IntVarP(&relistFlags.price, "price", "p", 0, "Item price per unit")
	relistCmd.Flags().IntVarP(&relistFlags.quantity, "quantity", "q", 0, "Item quantity (0 for all)")
}
