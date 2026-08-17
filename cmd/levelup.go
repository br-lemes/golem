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

var levelupExcludedSkills = []string{
	"gearcrafting",
	"jewelrycrafting",
	"weaponcrafting",
}

var levelupCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "levelup [account]",
	Short: "Show skill level-up requirements for an account",
	Long: `Show skill level-up requirements for an account

The levels are grouped in ranges of ten. Characters below a range reached by
another character are listed with the level they need to reach.

Arguments:
  account   The name of the account (optional).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		account := ""
		if len(args) > 0 {
			account = args[0]
		}
		characters, err := api.AccountsCharacters(account)
		if err != nil {
			return err
		}
		return console.Auto(skillLevelupRequirements(characters))
	},
}

func init() {
	rootCmd.AddCommand(levelupCmd)
}

func skillLevelupRequirements(characters []schemas.CharacterSchema) map[string]map[string]int {
	skills := slices.DeleteFunc(slices.Clone(database.Enum("CharacterLeaderboardType")), func(skill string) bool {
		return slices.Contains(levelupExcludedSkills, skill)
	})

	targets := map[string]int{}
	for _, skill := range skills {
		maxLevel := 0
		for _, character := range characters {
			level, exists := utils.GetCharacterSkillLevel(character, skill)
			if exists {
				maxLevel = max(maxLevel, level)
			}
		}
		targets[skill] = (maxLevel / 10) * 10
	}

	result := map[string]map[string]int{}
	for _, character := range characters {
		for skill, target := range targets {
			level, _ := utils.GetCharacterSkillLevel(character, skill)
			if level >= target || target == 0 {
				continue
			}
			if result[character.Name] == nil {
				result[character.Name] = map[string]int{}
			}
			result[character.Name][skill] = level
		}
	}
	return result
}
