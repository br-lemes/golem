package cmd

import (
	"fmt"

	. "github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress [account]",
	Short: "Show in-progress achievements for an account",
	Long: `Show in-progress achievements for an account
Arguments:
  account   The name of the account.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		page := 1
		result := []AccountAchievementSchema{}
		for {
			achievements, err := apiAccountAchievements(args[0], page)
			if err != nil {
				return err
			}
			for _, achievement := range achievements.Data {
				if len(achievement.Objectives) == 1 {
					progress := *achievement.Objectives[0].Progress
					total := achievement.Objectives[0].Total
					if progress == 0 || progress == total {
						continue
					}
					result = append(result, achievement)
					continue
				}
				for _, objective := range achievement.Objectives {
					if *objective.Progress > 0 {
						result = append(result, achievement)
						break
					}
				}
			}
			if page >= *achievements.Pages {
				break
			}
			page++
		}
		output(progressFormat(result))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(progressCmd)
}

func progressFormat(achievements []AccountAchievementSchema) map[string][]string {
	result := map[string][]string{}
	for _, achievement := range achievements {
		name := fmt.Sprintf("%s (%s)", achievement.Name, achievement.Code)
		for _, objective := range achievement.Objectives {
			result[name] = append(result[name], fmt.Sprintf("%s (%d/%d)",
				objective.Type, *objective.Progress, objective.Total))
		}
	}
	return result
}
