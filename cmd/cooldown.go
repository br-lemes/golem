package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var cooldownCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "cooldown [route]",
	Short: "List API routes that trigger a character cooldown",
	Long: `List API routes that trigger a character cooldown

Arguments:
  route   Path of a specific route to check for cooldown.`,
	ValidArgsFunction: completion.Custom(1, utils.GetRoutesCompletion).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			cooldown, err := utils.GetCooldown(args[0])
			if err != nil {
				return err
			}
			return console.Auto(map[string]bool{args[0]: cooldown})
		}

		cooldowns, err := utils.GetCooldowns()
		if err != nil {
			return err
		}

		return console.Auto(cooldowns)
	},
}

func init() {
	rootCmd.AddCommand(cooldownCmd)
}
