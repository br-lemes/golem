package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var levelGroups = []string{"skill", "character"}

type levelFlags struct {
	Group string   `flag:"group" shorthand:"g" desc:"Group by \"skill\" or \"character\" (defaults to \"character\" with --name, else \"skill\")"`
	Skill []string `flag:"skill" shorthand:"k" desc:"Show the level of specific skills"`
	Name  []string `flag:"name" shorthand:"n" desc:"Show the level of specific characters"`
}

var levelOptions levelFlags

var levelCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "level [account]",
	Short: "Show the level of an account",
	Long: `Show the level of an account

Arguments:
  account   The name of the account (optional).`,
	ValidArgsFunction: completion.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		levelOptions, err = utils.ReadFlags[levelFlags](cmd)
		if err != nil {
			return err
		}
		groupChanged := cmd.Flags().Changed("group")
		if !groupChanged {
			if len(levelOptions.Name) > 0 {
				levelOptions.Group = "character"
			} else {
				levelOptions.Group = "skill"
			}
		}
		if !slices.Contains(levelGroups, levelOptions.Group) {
			return fmt.Errorf("invalid group %q: allowed values are %v", levelOptions.Group, levelGroups)
		}
		validSkills := database.Enum("CharacterLeaderboardType")
		for _, skill := range levelOptions.Skill {
			if !slices.Contains(validSkills, skill) {
				return fmt.Errorf("invalid skill %q: allowed values are %v", skill, validSkills)
			}
		}
		validCharacters := utils.GetCharacters()
		for _, name := range levelOptions.Name {
			if !slices.Contains(validCharacters, name) {
				return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		account := ""
		if len(args) > 0 {
			account = args[0]
		}
		characters, err := api.AccountsCharacters(account)
		if err != nil {
			return err
		}
		if len(levelOptions.Name) > 0 {
			characters = slices.DeleteFunc(characters, func(c schemas.CharacterSchema) bool {
				return !slices.Contains(levelOptions.Name, c.Name)
			})
		}
		switch levelOptions.Group {
		case "character":
			return groupByCharacter(characters, levelOptions.Skill)
		case "skill":
			return groupBySkill(characters, levelOptions.Skill)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(levelCmd)
	err := utils.RegisterFlags[levelFlags](levelCmd)
	if err != nil {
		panic(err)
	}
	err = levelCmd.RegisterFlagCompletionFunc("group", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return levelGroups, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = levelCmd.RegisterFlagCompletionFunc("skill", completion.StringSlice(func() []string {
		return database.Enum("CharacterLeaderboardType")
	}))
	if err != nil {
		panic(err)
	}
	err = levelCmd.RegisterFlagCompletionFunc("name", completion.StringSlice(utils.GetCharacters))
	if err != nil {
		panic(err)
	}
}

func groupBySkill(characters []schemas.CharacterSchema, filterSkills []string) error {
	skills := database.Enum("CharacterLeaderboardType")
	if len(filterSkills) > 0 {
		skills = slices.DeleteFunc(skills, func(s string) bool {
			return !slices.Contains(filterSkills, s)
		})
	}
	levels := map[string]map[string]int{}
	for _, skill := range skills {
		levels[skill] = map[string]int{}
		for _, character := range characters {
			levels[skill][character.Name], _ = utils.GetCharacterSkillLevel(character, skill)
		}
	}
	return console.Auto(levels)
}

func groupByCharacter(characters []schemas.CharacterSchema, filterSkills []string) error {
	skills := database.Enum("CharacterLeaderboardType")
	if len(filterSkills) > 0 {
		skills = slices.DeleteFunc(skills, func(s string) bool {
			return !slices.Contains(filterSkills, s)
		})
	}
	levels := map[string]map[string]int{}
	for _, character := range characters {
		levels[character.Name] = map[string]int{}
		for _, skill := range skills {
			levels[character.Name][skill], _ = utils.GetCharacterSkillLevel(character, skill)
		}
	}
	return console.Auto(levels)
}
