package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type fightFlags struct {
	Food     string `flag:"food" desc:"auto-stock food from bank"`
	FoodOnly string `flag:"food-only" desc:"use only this food code"`
	NoFood   bool   `flag:"no-food" desc:"do not use food"`
	Utility1 string `flag:"utility1" desc:"item code to auto-refill in utility1 slot"`
	Utility2 string `flag:"utility2" desc:"item code to auto-refill in utility2 slot"`
}

var fightCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "fight <name> <code>",
	Short: "Fight continuously",
	Long: `Fight continuously

Arguments:
  name   Name of your character.
  code   The code of the monster.`,
	ValidArgsFunction: completion.CharacterName(1).Monster(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		flags, err := utils.ReadFlags[fightFlags](cmd)
		if err != nil {
			return err
		}
		err = fightValidate(code, flags, false)
		if err != nil {
			return err
		}
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		err = fightFoodValidate(character, flags)
		if err != nil {
			return err
		}
		if taskCompleted(character, code) {
			return nil
		}

		monster, _ := database.Monsters.Get(code)
		fightResult, err := best.FindFightByName(name, *monster, false, false)
		if err != nil {
			return err
		}
		if fightResult.Winrate < 100 {
			return fmt.Errorf("cannot safely fight %s: simulated winrate is %.2f%%", monster.Code, fightResult.Winrate)
		}
		if !fightUtilitiesMatch(fightResult, flags.Utility1, flags.Utility2) {
			return fmt.Errorf("fight utilities do not match the simulated configuration")
		}
		equipments := make([]schemas.EquipSchema, 0, len(fightResult.Equipment))
		for slot, code := range fightResult.Equipment {
			equipments = append(equipments, schemas.EquipSchema{
				Code: code,
				Slot: schemas.ItemSlot(slot),
			})
		}
		character, err = routine.Equip(name, equipments)
		if err != nil {
			return err
		}
		for {
			err = prepare(character, *monster, flags)
			if err != nil {
				return err
			}

			fightResult, err := api.MyActionFight(name, []string{})
			if err != nil {
				return err
			}
			character = fightResult.Characters[0]

			if taskCompleted(character, code) {
				return nil
			}
		}
	},
}

func fightValidate(monster string, flags fightFlags, boss bool) error {
	var found bool
	if boss {
		_, found = database.Bosses.Get(monster)
	} else {
		_, found = database.Monsters.Get(monster)
	}
	if !found {
		return fmt.Errorf("monster %s not found", monster)
	}
	if flags.NoFood && (flags.Food != "" || flags.FoodOnly != "") {
		return fmt.Errorf("--no-food cannot be combined with --food or --food-only")
	}
	if flags.Food != "" && flags.FoodOnly != "" {
		return fmt.Errorf("--food and --food-only cannot be combined")
	}
	if flags.NoFood {
		return nil
	}
	food := flags.Food
	if flags.FoodOnly != "" {
		food = flags.FoodOnly
	}
	if food == "" || food == "auto" {
		return nil
	}
	item, found := database.Items().Foods().Get(food)
	if !found {
		return fmt.Errorf("food %s not found", food)
	}
	_ = item
	return nil
}

func fightFoodValidate(character schemas.CharacterSchema, flags fightFlags) error {
	if flags.NoFood {
		return nil
	}
	food := flags.Food
	if flags.FoodOnly != "" {
		food = flags.FoodOnly
	}
	if food == "" || food == "auto" {
		return nil
	}
	item, _ := database.Items().Foods().Get(food)
	if !utils.MeetsItemConditions(character, *item) {
		return fmt.Errorf("food %s conditions not met for character %s", food, character.Name)
	}
	return nil
}

func fightUtilitiesMatch(result best.Result, utility1, utility2 string) bool {
	return result.Utilities["utility1"] == utility1 && result.Utilities["utility2"] == utility2
}

func prepare(character schemas.CharacterSchema, monster schemas.MonsterSchema, flags fightFlags) error {
	routine.Cooldown(character)
	character, err := routine.Inventory(character, []string{"food"})
	if err != nil {
		return err
	}
	character, err = routine.Bank(character, routine.BankOptions{
		Food:     flags.Food,
		FoodOnly: flags.FoodOnly != "",
		NoFood:   flags.NoFood,
		Utility1: flags.Utility1,
		Utility2: flags.Utility2,
	})
	if err != nil {
		return err
	}
	minHp := min(monster.Hp+(monster.Hp*20/100), character.MaxHp)
	if healingPotion(flags.Utility1) || healingPotion(flags.Utility2) {
		minHp = character.MaxHp
	}
	character, err = routine.Hp(character, routine.HpOptions{
		MinHP:    minHp,
		UseFood:  !flags.NoFood,
		FoodOnly: flags.FoodOnly,
	})
	if err != nil {
		return err
	}
	_, err = routine.Move(character, monster.Code)
	if err != nil {
		return err
	}
	return nil
}

func healingPotion(code string) bool {
	if code == "" {
		return false
	}
	potion, found := database.Items().Potions().Get(code)
	if !found || potion.Effects == nil {
		return false
	}
	for _, effect := range *potion.Effects {
		if effect.Code == "restore" {
			return true
		}
	}
	return false
}

func taskCompleted(character schemas.CharacterSchema, code string) bool {
	if character.TaskType == "monsters" && character.Task == code && character.TaskProgress == character.TaskTotal {
		console.Printf("  Task completed\n")
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(fightCmd)
	err := utils.RegisterFlags[fightFlags](fightCmd)
	if err != nil {
		panic(err)
	}
	fightCmd.Flags().Lookup("food").NoOptDefVal = "auto"
	err = fightCmd.RegisterFlagCompletionFunc("food", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return append([]string{"auto"}, database.Items().Foods().Keys()...), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = fightCmd.RegisterFlagCompletionFunc("food-only", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return database.Items().Foods().Keys(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}
