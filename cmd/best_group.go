package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	bestgroup "github.com/br-lemes/golem/pkg/best/group"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var bestGroupCmd = &cobra.Command{
	Args:              cobra.RangeArgs(2, 4),
	Use:               "group <monster> <name> [name...]",
	Short:             "Find equipment for a combat party",
	ValidArgsFunction: completion.Monster(1).CharacterName(3).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		monster, err := bestGroupMonster(args[0])
		if err != nil {
			return err
		}
		party, err := bestGroupParty(args[1:])
		if err != nil {
			return err
		}
		bank, err := api.MyBankItems()
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		result, err := bestgroup.FindFightWithAvailable(party, *monster, bestgroup.AvailableItems(bank, party))
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func bestGroupMonster(code string) (*schemas.MonsterSchema, error) {
	monster, found := catalog.Monsters.Get(code)
	if !found {
		return nil, fmt.Errorf("invalid monster: %s", code)
	}
	return monster, nil
}

func bestGroupParty(names []string) ([]schemas.CharacterSchema, error) {
	seen := map[string]bool{}
	party := make([]schemas.CharacterSchema, 0, len(names))
	for _, name := range names {
		if seen[name] {
			return nil, fmt.Errorf("duplicate character: %s", name)
		}
		seen[name] = true
		character, err := api.Characters(name)
		if err != nil {
			return nil, err
		}
		party = append(party, character)
	}
	return party, nil
}

func init() {
	bestCmd.AddCommand(bestGroupCmd)
}
