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
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var relistData struct {
	character schemas.CharacterSchema
	orders    []schemas.GEOrderSchema
}

type relistFlags struct {
	Price    int `flag:"price" shorthand:"p" desc:"Item price per unit"`
	Quantity int `flag:"quantity" shorthand:"q" desc:"Item quantity (0 for all)"`
}

var relistOptions relistFlags

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

		flags, err := utils.ReadFlags[relistFlags](cmd)
		if err != nil {
			return err
		}
		relistOptions = flags

		validCharacters := cache.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
		}
		_, found := database.Items().Tradeables().Get(code)
		if !found {
			return fmt.Errorf("item %q not tradeable or not found", code)
		}
		if relistOptions.Price <= 0 {
			return fmt.Errorf("price must be greater than 0")
		}
		if relistOptions.Quantity < 0 {
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
		if relistOptions.Quantity > 0 {
			totalToRelist = relistOptions.Quantity
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
				Price:    relistOptions.Price,
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
	err := utils.RegisterFlags[relistFlags](relistCmd)
	if err != nil {
		panic(err)
	}
}
