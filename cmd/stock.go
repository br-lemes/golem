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

type stockFlags struct {
	Target bool `flag:"target" desc:"Filter items below the target quantity instead of safety"`
}

var stockCmd = &cobra.Command{
	Use:   "stock",
	Short: "Stock",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[stockFlags](cmd)
		if err != nil {
			return err
		}
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		taskCodes := database.Tasks.Keys()
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
		skills := database.TaskSkills()
		maxSkillLevel := map[string]int{}
		for _, character := range characters {
			for _, skill := range skills {
				level, _ := utils.GetCharacterSkillLevel(character, skill)
				maxSkillLevel[skill] = max(maxSkillLevel[skill], level)
			}
		}
		result := map[string]ItemStock{}
		for _, task := range database.Tasks.All() {
			if task.Skill == nil || task.Type != "items" || task.Level > maxSkillLevel[*task.Skill] || slices.Contains(forbiddenTaskItem, task.Code) {
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
			target := safety + safety/3
			threshold := safety
			if flags.Target {
				threshold = target
			}
			if itemQuantity >= threshold {
				continue
			}
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
	err := utils.RegisterFlags[stockFlags](stockCmd)
	if err != nil {
		panic(err)
	}
}
