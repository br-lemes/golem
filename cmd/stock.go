package cmd

import (
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
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
	Target         bool `flag:"target" desc:"Filter items below the target quantity instead of safety"`
	IncludeLevel50 bool `flag:"include-level-50" desc:"Include level 50 tasks"`
}

var stockCmd = &cobra.Command{
	Use:   "stock [code...]",
	Short: "Show stock requirements for tasks and potions",
	Long: `Show stock requirements for tasks and potions

Arguments:
  code   The code of the task or potion.`,
	ValidArgsFunction: completion.Custom(0, stockItemCodes).Build(),
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
		taskCodes := stockItemCodes()
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
		skills := database.Tasks().Skills()
		maxSkillLevel := map[string]int{}
		for _, character := range characters {
			for _, skill := range skills {
				level, _ := utils.GetCharacterSkillLevel(character, skill)
				maxSkillLevel[skill] = max(maxSkillLevel[skill], level)
			}
		}
		result := map[string]ItemStock{}
		for _, task := range database.Tasks().All() {
			if !includeStockTask(task, args, flags, maxSkillLevel) {
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
			addStockResult(result, task.Code, itemQuantity, safety, target)
		}
		for _, potion := range database.Items().Potions().All() {
			if !includeStockPotion(potion, args, flags) {
				continue
			}
			users := 0
			for _, character := range characters {
				if canStockPotionForCharacter(character, potion) {
					users++
				}
			}
			if users == 0 {
				continue
			}
			safety := users * 100
			target := safety + safety/3
			itemQuantity := itemStockQuantity(itemsMap, potion.Code)
			addStockResult(result, potion.Code, itemQuantity, safety, target)
		}
		for code, stock := range result {
			threshold := stock.Safety
			if flags.Target {
				threshold = stock.Target
			}
			if stock.Current >= threshold {
				delete(result, code)
			}
		}
		return console.Auto(result)
	},
}

func stockItemCodes() []string {
	return append(database.Tasks().Items(), database.Items().Potions().Keys()...)
}

func itemStockQuantity(items map[string]*schemas.SimpleItemSchema, code string) int {
	item, exists := items[code]
	if !exists {
		return 0
	}
	return item.Quantity
}

func addStockResult(result map[string]ItemStock, code string, current, safety, target int) {
	stock := result[code]
	stock.Current = current
	stock.Safety += safety
	stock.Target += target
	stock.Needed = stock.Target - current
	result[code] = stock
}

func includeStockPotion(potion *schemas.ItemSchema, codes []string, flags stockFlags) bool {
	return (len(codes) == 0 || slices.Contains(codes, potion.Code)) && (flags.IncludeLevel50 || potion.Level != 50)
}

func canStockPotionForCharacter(character schemas.CharacterSchema, potion *schemas.ItemSchema) bool {
	if potion.Code == "small_health_potion" && (character.Level < 5 || character.Level > 19) {
		return false
	}
	return utils.MeetsItemConditions(character, *potion)
}

func includeStockTask(task *schemas.TaskFullSchema, codes []string, flags stockFlags, maxSkillLevel map[string]int) bool {
	if task.Skill == nil || task.Type != "items" || (!flags.IncludeLevel50 && task.Level == 50) {
		return false
	}
	return (len(codes) == 0 || slices.Contains(codes, task.Code)) && task.Level <= maxSkillLevel[*task.Skill] && !slices.Contains(forbiddenTaskItem, task.Code)
}

func init() {
	rootCmd.AddCommand(stockCmd)
	err := utils.RegisterFlags[stockFlags](stockCmd)
	if err != nil {
		panic(err)
	}
}
