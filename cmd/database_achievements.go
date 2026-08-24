package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbAchievementsCmd = &cobra.Command{
	Use:   "achievements",
	Short: "Achievements",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.Achievements(api.AchievementsOptions{})
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbAchievementsCmd)
}
