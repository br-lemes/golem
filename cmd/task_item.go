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
			if !cancelMissing && !isForbidden && bankQty < needed {
				return fmt.Errorf("missing item: %s (have %d, need %d)", character.Task, bankQty, needed)
			}
			if !isForbidden && bankQty >= needed {
				cancelsInARow = 0
				for character.TaskProgress < character.TaskTotal {
					character, err = routine.Move(character, "bank")
					if err != nil {
						return err
					}
					remaining := character.TaskTotal - character.TaskProgress
					toWithdraw := min(remaining, taskItemInvSpace(character, tasksCoinBuffer))
					withdraw := []schemas.SimpleItemSchema{
						{Code: character.Task, Quantity: toWithdraw},
					}
					coinsHeld := taskItemInvQty(character, tasksCoin)
					if coinsHeld < tasksCoinBuffer {
						coinsInBank := taskItemBankQty(bankItems, tasksCoin)
						topUp := min(tasksCoinBuffer-coinsHeld, coinsInBank)
						if topUp > 0 {
							withdraw = append(withdraw, schemas.SimpleItemSchema{
								Code:     tasksCoin,
								Quantity: topUp,
							})
						}
					}
					bankData, err := api.MyActionBankWithdrawItem(name, withdraw)
					if err != nil {
						return err
					}
					character = bankData.Character
					if coinsHeld > tasksCoinBuffer {
						excess := coinsHeld - tasksCoinBuffer
						bankData, err = api.MyActionBankDepositItem(name, []schemas.SimpleItemSchema{
							{Code: tasksCoin, Quantity: excess},
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
					trade, err := api.MyActionTaskTrade(name, schemas.SimpleItemSchema{
						Code:     character.Task,
						Quantity: toWithdraw,
					})
					if err != nil {
						return err
					}
					character = trade.Character
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

func taskItemInvSpace(character schemas.CharacterSchema, coinBuffer int) int {
	used := coinBuffer
	if character.Inventory != nil {
		for _, slot := range *character.Inventory {
			used += slot.Quantity
			if slot.Code != tasksCoin {
				used += slot.Quantity
			}
		}
	}
	return character.InventoryMaxItems - used
}
