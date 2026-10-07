package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type deleteFlags struct {
	All      bool `flag:"all" desc:"Delete all of this item from inventory and bank"`
	Discard  bool `flag:"discard" desc:"Delete items marked for discard from inventory and bank"`
	Quantity int  `flag:"quantity" shorthand:"q" desc:"Item quantity to delete"`
}

var deleteCmd = &cobra.Command{
	Args:  cobra.RangeArgs(1, 2),
	Use:   "delete <name> [code]",
	Short: "Delete items from a character's inventory and bank",
	Long: `Delete items from a character's inventory and bank

Arguments:
  name   Name of your character.
  code   The code of the item, unless --discard is used.`,
	ValidArgsFunction: completion.CharacterName(1).Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[deleteFlags](cmd)
		if err != nil {
			return err
		}
		err = deleteValidate(args, flags, cmd.Flags().Changed("quantity"))
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		movement, err := movementOptions(cmd)
		if err != nil {
			return err
		}
		if flags.Discard {
			return deleteDiscardRun(args[0], movement)
		}
		return deleteRun(args[0], args[1], flags, movement)
	},
}

func deleteValidate(args []string, options deleteFlags, quantityChanged bool) error {
	if options.Discard {
		if len(args) != 1 {
			return fmt.Errorf("--discard does not accept an item code")
		}
		if options.All || quantityChanged {
			return fmt.Errorf("--discard cannot be combined with --all or --quantity")
		}
		return nil
	}
	if len(args) != 2 {
		return fmt.Errorf("item code is required unless --discard is used")
	}
	if options.All && quantityChanged {
		return fmt.Errorf("--all cannot be combined with --quantity")
	}
	if !options.All && !quantityChanged {
		return fmt.Errorf("one of --quantity or --all is required")
	}
	if quantityChanged && options.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	_, found := catalog.Items().Get(args[1])
	if !found {
		return fmt.Errorf("item %q not found", args[1])
	}
	return nil
}

func deleteDiscardRun(name string, movement routine.MoveOptions) error {
	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	bankItems, err := api.MyBankItems()
	if err != nil {
		return err
	}
	quantities := map[string]int{}
	if character.Inventory != nil {
		for _, item := range *character.Inventory {
			quantities[item.Code] += item.Quantity
		}
	}
	for _, item := range bankItems {
		quantities[item.Code] += item.Quantity
	}
	for _, code := range cache.DiscardStockCodes() {
		if quantities[code] == 0 {
			continue
		}
		err = deleteRun(name, code, deleteFlags{All: true}, movement)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteRun(name, code string, options deleteFlags, movement routine.MoveOptions) error {
	character, err := api.Characters(name)
	if err != nil {
		return err
	}
	bankItems, err := api.MyBankItems()
	if err != nil {
		return err
	}
	bankQuantity := 0
	for _, item := range bankItems {
		if item.Code == code {
			bankQuantity += item.Quantity
		}
	}
	inventoryQuantity := 0
	if character.Inventory != nil {
		for _, item := range *character.Inventory {
			if item.Code == code {
				inventoryQuantity += item.Quantity
			}
		}
	}
	available := inventoryQuantity + bankQuantity
	if options.All {
		options.Quantity = available
		if available == 0 {
			return fmt.Errorf("no %q items available to delete", code)
		}
	}
	if options.Quantity > available {
		return fmt.Errorf("not enough items: required %d, available %d", options.Quantity, available)
	}

	routine.Cooldown(character)
	remaining := options.Quantity
	for remaining > 0 {
		if inventoryQuantity == 0 {
			character, err = routine.Move(character, "bank", movement)
			if err != nil {
				return err
			}
			items := routine.GetInventoryItems(character, routine.InventoryItemsOptions{
				KeepTravelPotions: true,
			})
			if len(items) > 0 {
				depositData, err := api.MyActionBankDepositItem(character.Name, items)
				if err != nil {
					return err
				}
				character = depositData.Character
			}
			character, err = routine.RestockTravelPotion(character, movement)
			if err != nil {
				return err
			}
			space := routine.SpaceAfterDeposit(character)
			withdrawQuantity := min(remaining, bankQuantity, space)
			if withdrawQuantity <= 0 {
				return fmt.Errorf("no inventory space available to withdraw %q", code)
			}
			withdrawData, err := api.MyActionBankWithdrawItem(character.Name, []schemas.SimpleItemSchema{{
				Code:     code,
				Quantity: withdrawQuantity,
			}})
			if err != nil {
				return err
			}
			character = withdrawData.Character
			inventoryQuantity += withdrawQuantity
			bankQuantity -= withdrawQuantity
		}

		deleteQuantity := min(remaining, inventoryQuantity)
		deleteData, err := api.MyActionDelete(character.Name, schemas.SimpleItemSchema{
			Code:     code,
			Quantity: deleteQuantity,
		})
		if err != nil {
			return err
		}
		character = deleteData.Character
		deleted := deleteData.Item.Quantity
		if deleted <= 0 {
			return fmt.Errorf("delete action removed no %q items", code)
		}
		inventoryQuantity -= deleted
		remaining -= deleted
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	err := utils.RegisterFlags[deleteFlags](deleteCmd)
	if err != nil {
		panic(err)
	}
}
