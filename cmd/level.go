package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var levelGroups = []string{"skill", "character"}

var levelFlags struct {
	group     string
	skill     []string
	character []string
}

var levelCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "level [account]",
	Short: "Show the level of an account",
	Long: `Show the level of an account

Arguments:
  account   The name of the account (optional).`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		groupChanged := cmd.Flags().Changed("group")
		if !groupChanged {
			if len(levelFlags.character) > 0 {
				levelFlags.group = "character"
			} else {
				levelFlags.group = "skill"
			}
		}
		if !slices.Contains(levelGroups, levelFlags.group) {
			return fmt.Errorf("invalid group %q: allowed values are %v", levelFlags.group, levelGroups)
		}
		validSkills := database.GetEnum("CharacterLeaderboardType")
		for _, skill := range levelFlags.skill {
			if !slices.Contains(validSkills, skill) {
				return fmt.Errorf("invalid skill %q: allowed values are %v", skill, validSkills)
			}
		}
		validCharacters := utils.GetCharacters()
		for _, character := range levelFlags.character {
			if !slices.Contains(validCharacters, character) {
				return fmt.Errorf("invalid character %q: allowed values are %v", character, validCharacters)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		account := ""
		if len(args) > 0 {
			account = args[0]
		}
		characters, err := api.AccountsCharacters(account)
		if err != nil {
			return err
		}
		if len(levelFlags.character) > 0 {
			characters = slices.DeleteFunc(characters, func(c schemas.CharacterSchema) bool {
				return !slices.Contains(levelFlags.character, c.Name)
			})
		}
		switch levelFlags.group {
		case "character":
			return groupByCharacter(characters, levelFlags.skill)
		case "skill":
			return groupBySkill(characters, levelFlags.skill)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(levelCmd)
	levelCmd.Flags().StringVarP(&levelFlags.group, "group", "g", "", `Group by "skill" or "character" (defaults to "character" with --character, else "skill")`)
	levelCmd.Flags().StringSliceVarP(&levelFlags.skill, "skill", "k", []string{}, "Show the level of specific skills")
	levelCmd.Flags().StringSliceVarP(&levelFlags.character, "character", "c", []string{}, "Show the level of specific characters")
	err := levelCmd.RegisterFlagCompletionFunc("group", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return levelGroups, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = levelCmd.RegisterFlagCompletionFunc("skill", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		skills := database.GetEnum("CharacterLeaderboardType")
		idx := strings.LastIndex(toComplete, ",")
		if idx == -1 {
			return skills, cobra.ShellCompDirectiveNoFileComp
		}
		prefix := toComplete[:idx]
		last := toComplete[idx+1:]
		suggestions := make([]string, 0, len(skills))
		for _, skill := range skills {
			hasPrefix := strings.HasPrefix(skill, last)
			if hasPrefix {
				suggestion := prefix + "," + skill
				suggestions = append(suggestions, suggestion)
			}
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = levelCmd.RegisterFlagCompletionFunc("character", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		characters := utils.GetCharacters()
		idx := strings.LastIndex(toComplete, ",")
		if idx == -1 {
			return characters, cobra.ShellCompDirectiveNoFileComp
		}
		prefix := toComplete[:idx]
		last := toComplete[idx+1:]
		suggestions := make([]string, 0, len(characters))
		for _, character := range characters {
			hasPrefix := strings.HasPrefix(character, last)
			if hasPrefix {
				suggestion := prefix + "," + character
				suggestions = append(suggestions, suggestion)
			}
		}
		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}

func groupBySkill(characters []schemas.CharacterSchema, filterSkills []string) error {
	skills := database.GetEnum("CharacterLeaderboardType")
	if len(filterSkills) > 0 {
		skills = slices.DeleteFunc(skills, func(s string) bool {
			return !slices.Contains(filterSkills, s)
		})
	}
	levels := map[string]map[string]int{}
	for _, skill := range skills {
		levels[skill] = map[string]int{}
		for _, character := range characters {
			levels[skill][character.Name] = utils.GetCharacterSkillLevel(character, skill)
		}
	}
	return console.Auto(levels)
}

func groupByCharacter(characters []schemas.CharacterSchema, filterSkills []string) error {
	skills := database.GetEnum("CharacterLeaderboardType")
	if len(filterSkills) > 0 {
		skills = slices.DeleteFunc(skills, func(s string) bool {
			return !slices.Contains(filterSkills, s)
		})
	}
	levels := map[string]map[string]int{}
	for _, character := range characters {
		levels[character.Name] = map[string]int{}
		for _, skill := range skills {
			levels[character.Name][skill] = utils.GetCharacterSkillLevel(character, skill)
		}
	}
	return console.Auto(levels)
}
