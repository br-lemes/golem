package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var progressCmd = &cobra.Command{
	Use:   "progress [account]",
	Short: "Show in-progress achievements for an account",
	Long: `Show in-progress achievements for an account

Arguments:
  account   The name of the account (optional).`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		account := ""
		if len(args) > 0 {
			account = args[0]
		}
		result := []schemas.AccountAchievementSchema{}
		achievements, err := api.AccountsAchievements(account)
		if err != nil {
			return err
		}
		for _, achievement := range achievements {
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
		return console.Auto(progressFormat(result))
	},
}

func init() {
	rootCmd.AddCommand(progressCmd)
}

func progressFormat(achievements []schemas.AccountAchievementSchema) map[string][]map[schemas.AchievementType]string {
	result := map[string][]map[schemas.AchievementType]string{}
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
			result[name] = append(result[name],
				map[schemas.AchievementType]string{
					objective.Type: fmt.Sprintf("%s (%d/%d)",
						targetVal, progressVal, objective.Total),
				})
		}
	}
	return result
}
