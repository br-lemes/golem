package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/best"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type fightFlags struct {
	Food            string   `flag:"food" desc:"auto-stock food from bank"`
	FoodOnly        string   `flag:"food-only" desc:"use only this food code"`
	FoodReserve     int      `flag:"food-reserve" default:"70" desc:"percentage of the planned inventory space for food"`
	NoFood          bool     `flag:"no-food" desc:"do not use food"`
	AllowUnsafe     bool     `flag:"allow-unsafe" desc:"allow fights with simulated winrate below 100%"`
	UseUtilities    bool     `flag:"use-utilities" desc:"use utilities selected by the simulator"`
	Preview         bool     `flag:"preview" desc:"show the planned inventory without taking actions"`
	Utility1        []string `flag:"utility1" desc:"item code(s) to auto-refill in utility1 slot; repeat per character"`
	Utility1Reserve int      `flag:"utility1-reserve" desc:"percentage of planned inventory space reserved for utility1"`
	Utility2        []string `flag:"utility2" desc:"item code(s) to auto-refill in utility2 slot; repeat per character"`
	Utility2Reserve int      `flag:"utility2-reserve" desc:"percentage of planned inventory space reserved for utility2"`
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
		movement, err := movementOptions(cmd)
		if err != nil {
			return err
		}
		err = fightValidate(code, flags, false)
		if err != nil {
			return err
		}
		options, err := flags.options(0, 1)
		if err != nil {
			return err
		}
		err = fightOptionsValidate(options)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		name = character.Name
		err = fightFoodValidate(character, options)
		if err != nil {
			return err
		}
		if !flags.Preview && taskCompleted(character, code) {
			return nil
		}

		monster, _ := catalog.Monsters.Get(code)
		fightResult, err := best.FindFightByName(name, *monster, false, false)
		if err != nil {
			return err
		}
		if fightResult.Winrate < 100 && !flags.AllowUnsafe {
			return fmt.Errorf("cannot safely fight %s: simulated winrate is %.2f%%", monster.Code, fightResult.Winrate)
		}
		utility1 := fightResult.FinalEquipment["utility1"]
		utility2 := fightResult.FinalEquipment["utility2"]
		utilitiesRequired := utility1 != "" || utility2 != ""
		if utilitiesRequired && !flags.UseUtilities {
			return fmt.Errorf("simulator requires utilities; use --use-utilities to allow them")
		}
		if !utilitiesRequired && flags.UseUtilities {
			return fmt.Errorf("simulator does not require utilities; remove --use-utilities")
		}
		if utilitiesRequired && (options.Utility1 != "" || options.Utility2 != "") {
			return fmt.Errorf("--utility1 and --utility2 cannot override utilities selected by the simulator")
		}
		if utilitiesRequired {
			options.Utility1 = utility1
			options.Utility2 = utility2
		}
		if flags.Preview {
			plan, err := routine.PlanInventory(character, movement.NoTeleport, options.foodReserve(), options.Utility1Reserve, options.Utility2Reserve)
			if err != nil {
				return err
			}
			return console.Auto(fightPreview{
				Monster: monster.Code,
				Winrate: float64(fightResult.Winrate),
				Characters: []fightCharacterPreview{{
					Character: name,
					Food:      options.Food,
					Utility1:  options.Utility1,
					Utility2:  options.Utility2,
					Inventory: plan,
				}},
			})
		}
		equipments := make([]schemas.EquipSchema, 0, len(fightResult.Equipment))
		for slot, code := range fightResult.Equipment {
			equipments = append(equipments, schemas.EquipSchema{
				Code: code,
				Slot: slot,
			})
		}
		character, err = routine.Equip(name, equipments, routine.EquipOptions{
			Movement:  movement,
			Utilities: fightResult.Utilities,
		})
		if err != nil {
			return err
		}
		for {
			err = prepare(character, *monster, options, movement)
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

type fightPreview struct {
	Monster    string                  `json:"monster"`
	Winrate    float64                 `json:"winrate,omitempty"`
	Characters []fightCharacterPreview `json:"characters"`
}

type fightCharacterPreview struct {
	Character string                `json:"character"`
	Food      string                `json:"food"`
	Utility1  string                `json:"utility1,omitempty"`
	Utility2  string                `json:"utility2,omitempty"`
	Inventory routine.InventoryPlan `json:"inventory"`
}

type fightOptions struct {
	Food            string
	FoodOnly        string
	NoFood          bool
	AllowUnsafe     bool
	UseUtilities    bool
	Utility1        string
	Utility2        string
	FoodReserve     int
	Utility1Reserve int
	Utility2Reserve int
}

func (flags fightFlags) options(index, characterCount int) (fightOptions, error) {
	utility1, err := utilityForCharacter(flags.Utility1, index, characterCount)
	if err != nil {
		return fightOptions{}, fmt.Errorf("--utility1: %w", err)
	}
	utility2, err := utilityForCharacter(flags.Utility2, index, characterCount)
	if err != nil {
		return fightOptions{}, fmt.Errorf("--utility2: %w", err)
	}
	return fightOptions{
		Food:            flags.Food,
		FoodOnly:        flags.FoodOnly,
		NoFood:          flags.NoFood,
		AllowUnsafe:     flags.AllowUnsafe,
		UseUtilities:    flags.UseUtilities,
		Utility1:        utility1,
		Utility2:        utility2,
		FoodReserve:     flags.FoodReserve,
		Utility1Reserve: flags.Utility1Reserve,
		Utility2Reserve: flags.Utility2Reserve,
	}, nil
}

func (flags fightOptions) foodReserve() int {
	if flags.NoFood || flags.Food == "" {
		return 0
	}
	return flags.FoodReserve
}

func utilityForCharacter(values []string, index, characterCount int) (string, error) {
	switch len(values) {
	case 0:
		return "", nil
	case 1:
		return normalizeUtility(values[0]), nil
	case characterCount:
		return normalizeUtility(values[index]), nil
	default:
		return "", fmt.Errorf("requires either 1 value or exactly %d values", characterCount)
	}
}

func normalizeUtility(value string) string {
	if value == "none" {
		return ""
	}
	return value
}

func fightValidate(monster string, flags fightFlags, boss bool) error {
	var found bool
	if boss {
		_, found = catalog.Bosses.Get(monster)
	} else {
		_, found = catalog.Monsters.Get(monster)
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
	for _, utility := range append(flags.Utility1, flags.Utility2...) {
		if utility == "" || utility == "none" {
			continue
		}
		_, found := catalog.Items().Potions().Get(utility)
		if !found {
			return fmt.Errorf("potion %s not found", utility)
		}
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
	item, found := catalog.Items().Foods().Get(food)
	if !found {
		return fmt.Errorf("food %s not found", food)
	}
	_ = item
	return nil
}

func fightOptionsValidate(flags fightOptions) error {
	if flags.Utility1 != "" && flags.Utility1 == flags.Utility2 {
		return fmt.Errorf("--utility1 and --utility2 cannot use the same item")
	}
	err := routine.ValidateInventoryReserve(flags.foodReserve(), flags.Utility1Reserve, flags.Utility2Reserve)
	if err != nil {
		return err
	}
	return nil
}

func fightFoodValidate(character schemas.CharacterSchema, flags fightOptions) error {
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
	item, _ := catalog.Items().Foods().Get(food)
	if !utils.MeetsItemConditions(character, *item) {
		return fmt.Errorf("food %s conditions not met for character %s", food, character.Name)
	}
	return nil
}

func prepare(character schemas.CharacterSchema, monster schemas.MonsterSchema, flags fightOptions, movement routine.MoveOptions) error {
	routine.Cooldown(character)
	character, err := routine.Inventory(character, []string{"food"}, movement)
	if err != nil {
		return err
	}
	food := flags.Food
	if flags.FoodOnly != "" {
		food = flags.FoodOnly
	}
	character, err = routine.Bank(character, routine.BankOptions{
		Movement:        movement,
		Food:            food,
		FoodOnly:        flags.FoodOnly != "",
		NoFood:          flags.NoFood,
		Utility1:        flags.Utility1,
		Utility2:        flags.Utility2,
		FoodReserve:     flags.FoodReserve,
		Utility1Reserve: flags.Utility1Reserve,
		Utility2Reserve: flags.Utility2Reserve,
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
	_, err = routine.Move(character, monster.Code, movement)
	if err != nil {
		return err
	}
	return nil
}

func healingPotion(code string) bool {
	if code == "" {
		return false
	}
	potion, found := catalog.Items().Potions().Get(code)
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

func registerFightFlags(cmd *cobra.Command) {
	err := utils.RegisterFlags[fightFlags](cmd)
	if err != nil {
		panic(err)
	}
	cmd.Flags().Lookup("food").NoOptDefVal = "auto"

	simpleCompletions := []struct {
		flag   string
		values func() []string
	}{
		{
			flag: "food",
			values: func() []string {
				return append([]string{"auto"}, catalog.Items().Foods().Keys()...)
			},
		},
		{
			flag: "food-only",
			values: func() []string {
				return catalog.Items().Foods().Keys()
			},
		},
	}
	for _, flagCompletion := range simpleCompletions {
		err := cmd.RegisterFlagCompletionFunc(flagCompletion.flag, func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) ([]string, cobra.ShellCompDirective) {
			return flagCompletion.values(), cobra.ShellCompDirectiveNoFileComp
		})
		if err != nil {
			panic(err)
		}
	}

	sliceCompletions := []struct {
		flag   string
		values func() []string
	}{
		{
			flag: "utility1",
			values: func() []string {
				return append([]string{"none"}, catalog.Items().Potions().Keys()...)
			},
		},
		{
			flag: "utility2",
			values: func() []string {
				return append([]string{"none"}, catalog.Items().Potions().Keys()...)
			},
		},
	}
	for _, flagCompletion := range sliceCompletions {
		err := cmd.RegisterFlagCompletionFunc(flagCompletion.flag, completion.StringSlice(flagCompletion.values))
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(fightCmd)
	registerFightFlags(fightCmd)
}
