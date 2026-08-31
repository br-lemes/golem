package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

const tasksCoin = "tasks_coin"

var forbiddenTaskItem = []string{
	"enchanted_mushroom",
	"magic_wood",
	"magical_plank",
	"strange_ore",
	"strangold_bar",
}

var taskItemCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "item <name>",
	Short: "Complete item-based tasks continuously",
	Long: `Complete item-based tasks continuously

Arguments:
  name   Name of your character.`,
	ValidArgsFunction: completion.CharacterName(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		tasksCoinBuffer, _ := cmd.Flags().GetInt("coin-buffer")
		tasksMaxCancel, _ := cmd.Flags().GetInt("max-cancel")
		cancelMissing, _ := cmd.Flags().GetBool("cancel-missing")

		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		if character.TaskType != "" && character.TaskType != "items" {
			return fmt.Errorf("has another task type: %s %s", character.TaskType, character.Task)
		}
		routine.Cooldown(character)

		cancelsInARow := 0
		for {
			if character.Task == "" {
				character, err = routine.Move(character, "items")
				if err != nil {
					return err
				}
				newTask, err := api.MyActionTaskNew(name)
				if err != nil {
					return err
				}
				character = newTask.Character
			}
			if character.TaskType != "items" {
				return fmt.Errorf("has another task type: %s %s", character.TaskType, character.Task)
			}
			bankItems, err := api.MyBankItems()
			if err != nil {
				return err
			}
			isForbidden := slices.Contains(forbiddenTaskItem, character.Task)
			bankQty := taskItemBankQty(bankItems, character.Task)
			needed := character.TaskTotal - character.TaskProgress
			available := bankQty + taskItemInvQty(character, character.Task)
			if !cancelMissing && !isForbidden && available < needed {
				return fmt.Errorf("missing item: %s (have %d, need %d)", character.Task, available, needed)
			}
			if !isForbidden && available >= needed {
				cancelsInARow = 0
				for character.TaskProgress < character.TaskTotal {
					remaining := character.TaskTotal - character.TaskProgress
					taskQuantity := taskItemInvQty(character, character.Task)
					toTrade := min(remaining, taskQuantity)
					if taskQuantity == 0 {
						toTrade = min(remaining, max(0, character.InventoryMaxItems-tasksCoinBuffer))
						character, err = routine.Move(character, "bank")
						if err != nil {
							return err
						}
						if character.Inventory != nil && len(*character.Inventory) > 0 {
							items := make([]schemas.SimpleItemSchema, 0, len(*character.Inventory))
							for _, item := range *character.Inventory {
								if item.Code == "" || item.Quantity <= 0 {
									continue
								}
								items = append(items, schemas.SimpleItemSchema{
									Code:     item.Code,
									Quantity: item.Quantity,
								})
							}
							if len(items) > 0 {
								bankData, err := api.MyActionBankDepositItem(name, items)
								if err != nil {
									return err
								}
								character = bankData.Character
							}
						}
						bankItems, err = api.MyBankItems()
						if err != nil {
							return err
						}
						withdraw := []schemas.SimpleItemSchema{
							{Code: character.Task, Quantity: toTrade},
						}
						coinsInBank := taskItemBankQty(bankItems, tasksCoin)
						if tasksCoinBuffer > 0 && coinsInBank > 0 {
							withdraw = append(withdraw, schemas.SimpleItemSchema{
								Code:     tasksCoin,
								Quantity: min(tasksCoinBuffer, coinsInBank),
							})
						}
						bankData, err := api.MyActionBankWithdrawItem(name, withdraw)
						if err != nil {
							return err
						}
						character = bankData.Character
					}
					character, err = routine.Move(character, "items")
					if err != nil {
						return err
					}
					trade, err := api.MyActionTaskTrade(name, schemas.SimpleItemSchema{
						Code:     character.Task,
						Quantity: toTrade,
					})
					if err != nil {
						return err
					}
					character = trade.Character
				}
				character, err = routine.Move(character, "items")
				if err != nil {
					return err
				}
				reward, err := api.MyActionTaskComplete(name)
				if err != nil {
					return err
				}
				character = reward.Character
				continue
			}
			if cancelsInARow >= tasksMaxCancel {
				return fmt.Errorf("reached %d consecutive task cancellations", tasksMaxCancel)
			}
			if taskItemInvQty(character, tasksCoin) < 1 {
				character, err = routine.Move(character, "bank")
				if err != nil {
					return err
				}
				bankData, err := api.MyActionBankWithdrawItem(name, []schemas.SimpleItemSchema{
					{Code: tasksCoin, Quantity: tasksCoinBuffer},
				})
				if err != nil {
					return err
				}
				character = bankData.Character
			}
			character, err = routine.Move(character, "items")
			if err != nil {
				return err
			}
			cancel, err := api.MyActionTaskCancel(name)
			if err != nil {
				return err
			}
			character = cancel.Character
			cancelsInARow++
		}
	},
}

func init() {
	taskCmd.AddCommand(taskItemCmd)
	taskItemCmd.Flags().IntP("coin-buffer", "b", 3, "Buffer of coins to keep in inventory")
	taskItemCmd.Flags().IntP("max-cancel", "c", 3, "Maximum consecutive task cancellations")
	taskItemCmd.Flags().Bool("cancel-missing", false, "Cancel task if required items are missing in bank")
}

func taskItemBankQty(items []schemas.SimpleItemSchema, code string) int {
	for _, item := range items {
		if item.Code == code {
			return item.Quantity
		}
	}
	return 0
}

func taskItemInvQty(character schemas.CharacterSchema, code string) int {
	if character.Inventory == nil {
		return 0
	}
	total := 0
	for _, slot := range *character.Inventory {
		if slot.Code == code {
			total += slot.Quantity
		}
	}
	return total
}
