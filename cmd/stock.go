package cmd

import (
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
		skills := database.GetTasksSkills()
		maxSkillLevel := map[string]int{}
		for _, character := range characters {
			for _, skill := range skills {
				maxSkillLevel[skill] = max(maxSkillLevel[skill],
					utils.GetCharacterSkillLevel(character, skill))
			}
		}
		maxQuantities := map[string]int{}
		tasks := database.GetTasksList()
		for _, task := range tasks {
			if task.Skill == nil {
				continue
			}
			if task.Level <= maxSkillLevel[*task.Skill] {
				maxQuantities[task.Code] = task.MaxQuantity
			}
		}
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		itemsMap := map[string]*schemas.SimpleItemSchema{}
		for _, item := range bankItems {
			_, exists := maxQuantities[item.Code]
			if !exists {
				continue
			}
			itemsMap[item.Code] = &item
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
						Quantity: 0,
					}
				}
				itemsMap[slot.Code].Quantity += slot.Quantity
			}
		}
		result := map[string]ItemStock{}
		for _, task := range tasks {
			item, exists := itemsMap[task.Code]
			if !exists {
				continue
			}
			safety := maxQuantities[item.Code]
			if safety > 100 {
				safety *= 5
			} else {
				safety *= 10
			}
			if item.Quantity >= safety {
				continue
			}
			target := safety + safety/3
			result[item.Code] = ItemStock{
				Current: item.Quantity,
				Needed:  target - item.Quantity,
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
