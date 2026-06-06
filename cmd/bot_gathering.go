package cmd

import (
	"fmt"
	"time"

	"github.com/br-lemes/golem/pkg/database"
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

		if !confirmSkill(name, string(resource.Skill)) {
			cmd.SilenceUsage = true
			return fmt.Errorf("operation cancelled")
		}

		character, err := apiCharacters(name)
		if err != nil {
			return err
		}
		handleCooldown(character)

		for {
			character, err = handleInventory(character)
			if err != nil {
				return err
			}

			character, err = handleMap(character, code)
			if err != nil {
				return err
			}

			fmt.Fprintf(writer, "[%s] Gathering resources...\n",
				time.Now().Format("15:04:05"))
			skill, err := apiActionGathering(name)
			if err != nil {
				return err
			}
			dropsStr := ""
			for _, item := range skill.Details.Items {
				dropsStr += fmt.Sprintf("%dx %s, ", item.Quantity, item.Code)
			}
			fmt.Fprintf(writer, "[%s] Gathered %s\n",
				time.Now().Format("15:04:05"), dropsStr)

			character = skill.Character
		}
	},
}

func init() {
	botCmd.AddCommand(botGatheringCmd)
}
