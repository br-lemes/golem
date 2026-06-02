package cmd

import (
	. "github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var dbEffectsCmd = &cobra.Command{
	Use:   "effects",
	Short: "Effects",
	RunE: func(cmd *cobra.Command, args []string) error {
		page := 1
		result := []EffectSchema{}
		for {
			effects, err := apiEffects(page)
			if err != nil {
				return err
			}
			result = append(result, effects.Data...)
			if page >= *effects.Pages {
				break
			}
			page++
		}
		output(result)
		return nil
	},
}

func init() {
	databaseCmd.AddCommand(dbEffectsCmd)
}
