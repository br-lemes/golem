package cmd

import (
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type ItemStock struct {
	Current int `json:"current"`
	Needed  int `json:"needed"`
	Safety  int `json:"safety"`
	Target  int `json:"target"`
}

var stockCmd = &cobra.Command{
	Use:   "stock",
	Short: "Stock",
	RunE: func(cmd *cobra.Command, args []string) error {
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		taskCodes := database.GetTasksCodes()
		itemsMap := map[string]*schemas.SimpleItemSchema{}
		for _, item := range bankItems {
			if slices.Contains(taskCodes, item.Code) {
				itemsMap[item.Code] = &item
			}
		}
		for _, character := range characters {
			if character.Inventory == nil {
				continue
			}
			for _, slot := range *character.Inventory {
				if slot.Code == "" {
					continue
				}
				_, exists := itemsMap[slot.Code]
				if !exists {
					itemsMap[slot.Code] = &schemas.SimpleItemSchema{
						Code:     slot.Code,
						Quantity: slot.Quantity,
					}
				}
				itemsMap[slot.Code].Quantity += slot.Quantity
			}
		}
		skills := database.GetTasksSkills()
		maxSkillLevel := map[string]int{}
		for _, character := range characters {
			for _, skill := range skills {
				maxSkillLevel[skill] = max(maxSkillLevel[skill],
					utils.GetCharacterSkillLevel(character, skill))
			}
		}
		tasks := database.GetTasksList()
		result := map[string]ItemStock{}
		for _, task := range tasks {
			if task.Skill == nil ||
				task.Type != "items" ||
				task.Level > maxSkillLevel[*task.Skill] ||
				slices.Contains(forbiddenTaskItem, task.Code) {
				continue
			}
			var itemQuantity int
			item, exists := itemsMap[task.Code]
			if exists {
				itemQuantity = item.Quantity
			} else {
				itemQuantity = 0
			}
			safety := task.MaxQuantity
			if safety > 100 {
				safety *= 5
			} else {
				safety *= 10
			}
			if itemQuantity >= safety {
				continue
			}
			target := safety + safety/3
			result[task.Code] = ItemStock{
				Current: itemQuantity,
				Needed:  target - itemQuantity,
				Safety:  safety,
				Target:  target,
			}
		}
		return console.Auto(result)
	},
}

func init() {
	rootCmd.AddCommand(stockCmd)
}
