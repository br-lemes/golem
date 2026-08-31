package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
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

var levelCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "level [account]",
	Short: "Show the level of an account",
	Long: `Show the level of an account

Arguments:
  account   The name of the account (optional).`,
	ValidArgsFunction: completion.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		account := ""
		if len(args) > 0 {
			account = args[0]
		}
		options, err := utils.ReadFlags[levelFlags](cmd)
		if err != nil {
			return err
		}
		if !cmd.Flags().Changed("group") {
			if len(options.Name) > 0 {
				options.Group = "character"
			} else {
				options.Group = "skill"
			}
		}
		err = levelValidate(options)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return levelRun(account, options)
	},
}

func levelValidate(options levelFlags) error {
	if !slices.Contains(levelGroups, options.Group) {
		return fmt.Errorf("invalid group %q: allowed values are %v", options.Group, levelGroups)
	}
	validSkills := database.Enum("CharacterLeaderboardType")
	for _, skill := range options.Skill {
		if !slices.Contains(validSkills, skill) {
			return fmt.Errorf("invalid skill %q: allowed values are %v", skill, validSkills)
		}
	}
	return nil
}

func levelRun(account string, levelOptions levelFlags) error {
	characters, err := api.AccountsCharacters(account)
	if err != nil {
		return err
	}
	if len(levelOptions.Name) > 0 {
		availableCharacters := make(map[string]struct{}, len(characters))
		for _, character := range characters {
			availableCharacters[character.Name] = struct{}{}
		}
		for _, name := range levelOptions.Name {
			_, found := availableCharacters[name]
			if !found {
				return fmt.Errorf("character %q not found in account %q", name, account)
			}
		}
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
}

func groupBySkill(characters []schemas.CharacterSchema, filterSkills []string) error {
	return console.Auto(levelsBySkill(characters, filterSkills))
}

func groupByCharacter(characters []schemas.CharacterSchema, filterSkills []string) error {
	bySkill := levelsBySkill(characters, filterSkills)
	levels := map[string]map[string]int{}
	for skill, characterLevels := range bySkill {
		for character, level := range characterLevels {
			if levels[character] == nil {
				levels[character] = map[string]int{}
			}
			levels[character][skill] = level
		}
	}
	return console.Auto(levels)
}

func levelsBySkill(characters []schemas.CharacterSchema, filterSkills []string) map[string]map[string]int {
	skills := slices.Clone(database.Enum("CharacterLeaderboardType"))
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
	return levels
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
	err = levelCmd.RegisterFlagCompletionFunc("name", completion.StringSlice(cache.GetCharacters))
	if err != nil {
		panic(err)
	}
}
