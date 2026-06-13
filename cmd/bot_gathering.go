package cmd

import (
	"fmt"
	"time"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/task"
	"github.com/spf13/cobra"
)

var botGatheringCmd = &cobra.Command{
	Use:   "gathering <name> <code>",
	Short: "Gather resources continuously",
	Long: `Gather resources continuously

Arguments:
  name   Name of your character.
  code   The code of the resource.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		resource, found := database.GetResource(code)
		if !found {
			return fmt.Errorf("resource %s not found", code)
		}

		if !config.ConfirmSkill(name, string(resource.Skill)) {
			cmd.SilenceUsage = true
			return fmt.Errorf("operation cancelled")
		}

		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		task.Cooldown(character)

		for {
			character, err = task.Inventory(character, false)
			if err != nil {
				return err
			}

			character, err = task.Move(character, code)
			if err != nil {
				return err
			}

			console.Printf("[%s] Gathering resources...\n",
				time.Now().Format("15:04:05"))
			skill, err := api.MyActionGathering(name)
			if err != nil {
				return err
			}
			dropsStr := ""
			for _, item := range skill.Details.Items {
				dropsStr += fmt.Sprintf("%dx %s, ", item.Quantity, item.Code)
			}
			console.Printf("[%s] Gathered %s\n",
				time.Now().Format("15:04:05"), dropsStr)

			character = skill.Character
		}
	},
}

func init() {
	botCmd.AddCommand(botGatheringCmd)
}
