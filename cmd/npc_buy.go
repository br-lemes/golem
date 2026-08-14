package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var npcBuyData struct {
	character     schemas.CharacterSchema
	bankGold      int
	bankItem      schemas.SimpleItemSchema
	inventoryItem schemas.SimpleItemSchema
	item          *schemas.NPCItemSchema
}

var npcBuyFlags struct {
	quantity int
}

var npcBuyCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "buy <name> <code>",
	Short: "buy",
	Long: `buy

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).NPCBuy(1).Build(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]
		validCharacters := utils.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
		}
		var exists bool
		npcBuyData.item, exists = database.NpcsItems.Get(code)
		if npcBuyData.item == nil || !exists {
			return fmt.Errorf("invalid item %q: allowed values are %v", args[1], completion.GetNPCBuyItems())
		}
		if npcBuyData.item.BuyPrice == nil || *npcBuyData.item.BuyPrice <= 0 {
			return fmt.Errorf("item %q does not have a buy price", code)
		}
		if npcBuyFlags.quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
		var err error
		npcBuyData.character, err = api.Characters(name)
		if err != nil {
			return err
		}
		if npcBuyData.item.Currency != "gold" && *npcBuyData.item.BuyPrice > npcBuyData.character.InventoryMaxItems {
			return fmt.Errorf("not enough inventory space: required %d, available %d", *npcBuyData.item.BuyPrice, npcBuyData.character.InventoryMaxItems)
		}
		if npcBuyData.item.Currency == "gold" {
			bank, err := api.MyBank()
			if err != nil {
				return err
			}
			npcBuyData.bankGold = bank.Gold
			totalGold := npcBuyData.character.Gold + npcBuyData.bankGold
			requiredGold := npcBuyFlags.quantity * *npcBuyData.item.BuyPrice
			if totalGold < requiredGold {
				return fmt.Errorf("not enough gold: required %d, available %d", requiredGold, totalGold)
			}
			return nil
		}
		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.Code == npcBuyData.item.Currency {
				npcBuyData.bankItem = item
				break
			}
		}
		npcBuyData.inventoryItem = schemas.SimpleItemSchema{}
		for _, item := range inventoryItems() {
			if item.Code == npcBuyData.item.Currency {
				npcBuyData.inventoryItem = schemas.SimpleItemSchema{
					Code:     item.Code,
					Quantity: item.Quantity,
				}
				break
			}
		}
		totalAvailable := npcBuyData.bankItem.Quantity + npcBuyData.inventoryItem.Quantity
		totalNeeded := npcBuyFlags.quantity * *npcBuyData.item.BuyPrice
		if totalAvailable < totalNeeded {
			return fmt.Errorf("not enough items: required %d, available %d", totalNeeded, totalAvailable)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		totalBought := 0
		if npcBuyData.item.Currency == "gold" {
			cost := npcBuyFlags.quantity * *npcBuyData.item.BuyPrice
			if npcBuyData.character.Gold < cost {
				err := moveBank()
				if err != nil {
					return err
				}
				err = depositAll()
				if err != nil {
					return err
				}
				err = withdrawGold(cost - npcBuyData.character.Gold)
				if err != nil {
					return err
				}
			}
		}
		for totalBought < npcBuyFlags.quantity {
			remaining := npcBuyFlags.quantity - totalBought
			quantity := min(remaining, npcBuyData.character.InventoryMaxItems)
			if npcBuyData.item.Currency != "gold" {
				maxByCurrency := npcBuyData.character.InventoryMaxItems / *npcBuyData.item.BuyPrice
				quantity = min(quantity, maxByCurrency)
			}
			cost := quantity * *npcBuyData.item.BuyPrice
			if npcBuyData.item.Currency == "gold" {
				for _, invItem := range inventoryItems() {
					if invItem.Code != "" {
						err := moveBank()
						if err != nil {
							return err
						}
						err = depositAll()
						if err != nil {
							return err
						}
						break
					}
				}
			} else {
				needBank := false
				for _, invItem := range inventoryItems() {
					if invItem.Code != "" && invItem.Code != npcBuyData.item.Currency {
						needBank = true
						break
					}
				}
				if needBank || npcBuyData.inventoryItem.Quantity < cost {
					err := moveBank()
					if err != nil {
						return err
					}
					err = depositAll()
					if err != nil {
						return err
					}
					err = withdrawItem(cost)
					if err != nil {
						return err
					}
				}
			}
			err := moveNpc()
			if err != nil {
				return err
			}
			err = buyItem(quantity)
			if err != nil {
				return err
			}
			totalBought += quantity
		}
		return nil
	},
}

func init() {
	npcCmd.AddCommand(npcBuyCmd)
	npcBuyCmd.Flags().IntVarP(&npcBuyFlags.quantity, "quantity", "q", 0, "Item quantity")
}

func inventoryItems() []schemas.InventorySlotSchema {
	if npcBuyData.character.Inventory == nil {
		return nil
	}
	return *npcBuyData.character.Inventory
}

func moveBank() error {
	var err error
	npcBuyData.character, err = routine.Move(npcBuyData.character, "bank")
	return err
}

func moveNpc() error {
	var err error
	npcBuyData.character, err = routine.Move(npcBuyData.character, npcBuyData.item.Npc)
	return err
}

func depositAll() error {
	name := npcBuyData.character.Name
	items := routine.GetInventoryItems(npcBuyData.character, nil)
	if len(items) == 0 {
		return nil
	}
	depositData, err := api.MyActionBankDepositItem(name, items)
	if err != nil {
		return err
	}
	npcBuyData.character = depositData.Character
	return nil
}

func withdrawGold(quantity int) error {
	name := npcBuyData.character.Name
	goldData, err := api.MyActionBankWithdrawGold(name, quantity)
	if err != nil {
		return err
	}
	npcBuyData.character = goldData.Character
	return nil
}

func withdrawItem(quantity int) error {
	name := npcBuyData.character.Name
	items := []schemas.SimpleItemSchema{
		{Code: npcBuyData.item.Currency, Quantity: quantity},
	}
	withdrawData, err := api.MyActionBankWithdrawItem(name, items)
	if err != nil {
		return err
	}
	npcBuyData.character = withdrawData.Character
	npcBuyData.inventoryItem.Quantity = quantity
	return nil
}

func buyItem(quantity int) error {
	name := npcBuyData.character.Name
	code := npcBuyData.item.Code
	buyItem := schemas.SimpleItemSchema{Code: code, Quantity: quantity}
	buyData, err := api.MyActionNPCBuy(name, buyItem)
	if err != nil {
		return err
	}
	npcBuyData.character = buyData.Character
	npcBuyData.inventoryItem.Quantity = 0
	return nil
}
