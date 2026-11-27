package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/account"
	"github.com/br-lemes/golem/pkg/best/damage"
	"github.com/br-lemes/golem/pkg/best/group"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type bestDamageFlags struct {
	All     bool     `flag:"all" desc:"Include all items from the catalog, even if not owned"`
	Weapon  string   `flag:"weapon" desc:"Show all combinations for a specific weapon"`
	Without []string `flag:"without" desc:"Exclude items from the selection"`
}

var bestDamageCmd = &cobra.Command{
	Args:  cobra.RangeArgs(2, 4),
	Use:   "damage <monster> <name> [name...]",
	Short: "Find the best damage gear against a monster",
	Long: `Find the best damage gear against a monster

Arguments:
  monster   The code of the monster.
  name      Name of your character.`,
	ValidArgsFunction: completion.Monster(1).CharacterName(3).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[bestDamageFlags](cmd)
		if err != nil {
			return err
		}
		err = bestDamageValidateWithout(flags.Without)
		if err != nil {
			return err
		}
		monster := args[0]
		names := args[1:]
		characters, err := bestDamageValidate(monster, names)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return bestDamageRun(monster, characters, flags)
	},
}

type bestDamageCharacter struct {
	Name    string          `json:"name"`
	Level   int             `json:"level"`
	Weapons []damage.Weapon `json:"weapons"`
}

type bestDamageResult struct {
	Monster    string                `json:"monster"`
	Characters []bestDamageCharacter `json:"characters"`
}

func bestDamageValidate(monster string, names []string) ([]schemas.CharacterSchema, error) {
	_, found := catalog.Monsters.Get(monster)
	if !found {
		return nil, fmt.Errorf("invalid monster: %s", monster)
	}
	_, isBoss := catalog.Bosses.Get(monster)
	if !isBoss && len(names) > 1 {
		return nil, fmt.Errorf("monster %s allows only 1 character", monster)
	}
	characters, err := account.Characters()
	if err != nil {
		return nil, err
	}
	available := make(map[string]schemas.CharacterSchema, len(characters))
	for _, character := range characters {
		available[character.Name] = character
	}
	selected := make([]schemas.CharacterSchema, 0, len(names))
	for _, name := range names {
		character, found := available[name]
		if !found {
			return nil, fmt.Errorf("character %q not found in account", name)
		}
		selected = append(selected, character)
	}
	if isBoss && len(names) == 1 {
		console.Warnf("monster %s is a boss and requires 2 to 3 characters\n", monster)
	}
	return selected, nil
}

func bestDamageValidateWithout(codes []string) error {
	for _, code := range codes {
		_, found := catalog.Items().Get(code)
		if !found {
			return fmt.Errorf("invalid item: %s", code)
		}
	}
	return nil
}

func bestDamageRun(monster string, characters []schemas.CharacterSchema, flags bestDamageFlags) error {
	options, err := bestDamageOptions(characters, flags)
	if err != nil {
		return err
	}
	monsterData, _ := catalog.Monsters.Get(monster)
	result := make([]bestDamageCharacter, 0, len(characters))
	for _, character := range characters {
		weapons, err := bestDamageWeapons(*monsterData, character, options, flags.Weapon)
		if err != nil {
			return err
		}
		result = append(result, bestDamageCharacter{
			Name:    character.Name,
			Level:   character.Level,
			Weapons: weapons,
		})
	}
	return console.Auto(bestDamageResult{Monster: monster, Characters: result})
}

func bestDamageOptions(characters []schemas.CharacterSchema, flags bestDamageFlags) (damage.Options, error) {
	options := damage.Options{Without: bestDamageWithout(flags.Without)}
	if flags.All {
		return options, nil
	}
	bank, err := account.BankItems()
	if err != nil {
		return damage.Options{}, err
	}
	options.Available = group.AvailableItems(bank, characters)
	return options, nil
}

func bestDamageWithout(codes []string) map[string]bool {
	if len(codes) == 0 {
		return nil
	}
	result := make(map[string]bool, len(codes))
	for _, code := range codes {
		result[code] = true
	}
	return result
}

func bestDamageWeapons(monster schemas.MonsterSchema, character schemas.CharacterSchema, options damage.Options, requested string) ([]damage.Weapon, error) {
	if requested == "" {
		return damage.FindWeapons(monster, character, options), nil
	}
	weapon, err := damage.FindWeapon(monster, character, options, requested)
	if err != nil {
		return nil, err
	}
	return []damage.Weapon{weapon}, nil
}

func init() {
	err := utils.RegisterFlags[bestDamageFlags](bestDamageCmd)
	if err != nil {
		panic(err)
	}
	err = bestDamageCmd.RegisterFlagCompletionFunc("weapon", completion.Weapon(0).Build())
	if err != nil {
		panic(err)
	}
	err = bestDamageCmd.RegisterFlagCompletionFunc("without", completion.UniqueList(catalog.Items().Equipments().Keys))
	if err != nil {
		panic(err)
	}
	bestCmd.AddCommand(bestDamageCmd)
}
