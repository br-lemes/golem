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
					if achievement.Objectives[0].Progress == nil {
						continue
					}
					progress := *achievement.Objectives[0].Progress
					total := achievement.Objectives[0].Total
					if progress == 0 || progress == total {
						continue
					}
					result = append(result, achievement)
					continue
				}
				hasProgress := false
				allCompleted := true
				for _, objective := range achievement.Objectives {
					if objective.Progress == nil {
						allCompleted = false
						continue
					}
					if *objective.Progress > 0 {
						hasProgress = true
					}
					if *objective.Progress < objective.Total {
						allCompleted = false
					}
				}
				if hasProgress && !allCompleted {
					result = append(result, achievement)
				}
			}
			if achievements.Pages == nil {
				break
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

func progressFormat(achievements []AccountAchievementSchema) map[string][]map[AchievementType]string {
	result := map[string][]map[AchievementType]string{}
	for _, achievement := range achievements {
		name := fmt.Sprintf("%s (%s)", achievement.Name, achievement.Code)
		for _, objective := range achievement.Objectives {
			targetVal := ""
			if objective.Target != nil {
				targetVal = *objective.Target
			}
			progressVal := 0
			if objective.Progress != nil {
				progressVal = *objective.Progress
			}
			result[name] = append(result[name], map[AchievementType]string{
				objective.Type: fmt.Sprintf("%s (%d/%d)",
					targetVal, progressVal, objective.Total),
			})
		}
	}
	return result
}
