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

var npcSellData struct {
	character     schemas.CharacterSchema
	bankItem      schemas.SimpleItemSchema
	inventoryItem schemas.SimpleItemSchema
	item          *schemas.NPCItemSchema
}

var npcSellFlags struct {
	quantity int
}

var npcSellCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "sell <name> <code>",
	Short: "sell",
	Long: `sell

Arguments:
  name   Name of your character.
  code   The code of the item.`,
	ValidArgsFunction: completion.CharacterName(1).NPCSell(1).Build(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]
		validCharacters := utils.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
		}
		var exists bool
		npcSellData.item, exists = database.NpcsItems.Get(code)
		if npcSellData.item == nil || !exists {
			return fmt.Errorf("invalid item %q: allowed values are %v", args[1], completion.GetNPCSellItems())
		}
		if npcSellData.item.SellPrice == nil || *npcSellData.item.SellPrice <= 0 {
			return fmt.Errorf("item %q does not have a sell price", code)
		}
		if npcSellFlags.quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
		var err error
		npcSellData.character, err = api.Characters(name)
		if err != nil {
			return err
		}
		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		for _, item := range items {
			if item.Code == npcSellData.item.Code {
				npcSellData.bankItem = item
				break
			}
		}
		npcSellData.inventoryItem = schemas.SimpleItemSchema{}
		for _, item := range npcSellInventoryItems() {
			if item.Code == npcSellData.item.Code {
				npcSellData.inventoryItem = schemas.SimpleItemSchema{
					Code:     item.Code,
					Quantity: item.Quantity,
				}
				break
			}
		}
		totalAvailable := npcSellData.bankItem.Quantity + npcSellData.inventoryItem.Quantity
		if totalAvailable < npcSellFlags.quantity {
			return fmt.Errorf("not enough items: required %d, available %d", npcSellFlags.quantity, totalAvailable)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		totalSold := 0
		for totalSold < npcSellFlags.quantity {
			remaining := npcSellFlags.quantity - totalSold
			quantity := min(remaining, npcSellData.character.InventoryMaxItems)
			needBank := false
			for _, invItem := range npcSellInventoryItems() {
				if invItem.Code != "" && invItem.Code != npcSellData.item.Code {
					needBank = true
					break
				}
			}
			if needBank || npcSellData.inventoryItem.Quantity < quantity {
				err := npcSellMoveBank()
				if err != nil {
					return err
				}
				err = npcSellDepositAll()
				if err != nil {
					return err
				}
				err = npcSellWithdrawItem(quantity)
				if err != nil {
					return err
				}
			}
			err := npcSellMoveNpc()
			if err != nil {
				return err
			}
			err = npcSellSellItem(quantity)
			if err != nil {
				return err
			}
			totalSold += quantity
		}
		return nil
	},
}

func init() {
	npcCmd.AddCommand(npcSellCmd)
	npcSellCmd.Flags().IntVarP(&npcSellFlags.quantity, "quantity", "q", 0, "Item quantity")
}

func npcSellInventoryItems() []schemas.InventorySlotSchema {
	if npcSellData.character.Inventory == nil {
		return nil
	}
	return *npcSellData.character.Inventory
}

func npcSellMoveBank() error {
	var err error
	npcSellData.character, err = routine.Move(npcSellData.character, "bank")
	return err
}

func npcSellMoveNpc() error {
	var err error
	npcSellData.character, err = routine.Move(npcSellData.character, npcSellData.item.Npc)
	return err
}

func npcSellDepositAll() error {
	name := npcSellData.character.Name
	items := routine.GetInventoryItems(npcSellData.character, nil)
	if len(items) == 0 {
		return nil
	}
	depositData, err := api.MyActionBankDepositItem(name, items)
	if err != nil {
		return err
	}
	npcSellData.character = depositData.Character
	return nil
}

func npcSellWithdrawItem(quantity int) error {
	name := npcSellData.character.Name
	items := []schemas.SimpleItemSchema{
		{Code: npcSellData.item.Code, Quantity: quantity},
	}
	withdrawData, err := api.MyActionBankWithdrawItem(name, items)
	if err != nil {
		return err
	}
	npcSellData.character = withdrawData.Character
	npcSellData.inventoryItem.Quantity = quantity
	return nil
}

func npcSellSellItem(quantity int) error {
	name := npcSellData.character.Name
	code := npcSellData.item.Code
	sellItem := schemas.SimpleItemSchema{Code: code, Quantity: quantity}
	sellData, err := api.MyActionNPCSell(name, sellItem)
	if err != nil {
		return err
	}
	npcSellData.character = sellData.Character
	npcSellData.inventoryItem.Quantity -= quantity
	return nil
}
