package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var levelCmd = &cobra.Command{
	Use:   "level [account]",
	Short: "Show the level of an account",
	Long: `Show the level of an account

Arguments:
  account   The name of the account (optional).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		account := ""
		if len(args) > 0 {
			account = args[0]
		}
		characters, err := api.AccountsCharacters(account)
		if err != nil {
			return err
		}
		group := cmd.Flags().Lookup("group").Value.String()
		switch group {
		case "character":
			groupByCharacter(characters)
		case "skill":
			groupBySkill(characters)
		}
		// skill := cmd.Flags().Lookup("skill").Value.String()
		// skills := database.GetEnum("CraftSkill")
		// levels := map[string]map[string]int{}
		// for _, character := range characters {
		// 	levels[character.Name] = map[string]int{}
		// 	for _, skill := range skills {
		// 		levels[character.Name][skill] =
		// 			utils.GetCharacterSkillLevel(character, skill)
		// 	}
		// }
		// console.Auto(levels)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(levelCmd)
	levelCmd.Flags().StringP("group", "g", "skill",
		"Group the output by skill or by character (default: skill)")
	// levelCmd.Flags().String("skill", "",
	// 	"Show the level of a specific skill (optional)")
	// levelCmd.RegisterFlagCompletionFunc("skill", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	return database.GetEnum("CraftSkill"),
	// 		cobra.ShellCompDirectiveNoFileComp
	// })
}

func groupBySkill(characters []schemas.CharacterSchema) {
	skills := database.GetEnum("CharacterLeaderboardType")
	levels := map[string]map[string]int{}
	for _, skill := range skills {
		levels[skill] = map[string]int{}
		for _, character := range characters {
			levels[skill][character.Name] =
				utils.GetCharacterSkillLevel(character, skill)
		}
	}
	console.Auto(levels)
}

func groupByCharacter(characters []schemas.CharacterSchema) {
	skills := database.GetEnum("CharacterLeaderboardType")
	levels := map[string]map[string]int{}
	for _, character := range characters {
		levels[character.Name] = map[string]int{}
		for _, skill := range skills {
			levels[character.Name][skill] =
				utils.GetCharacterSkillLevel(character, skill)
		}
	}
	console.Auto(levels)
}
