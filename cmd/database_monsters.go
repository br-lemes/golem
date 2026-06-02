package cmd

import (
	. "github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var dbMonstersCmd = &cobra.Command{
	Use:   "monsters",
	Short: "Monsters",
	RunE: func(cmd *cobra.Command, args []string) error {
		page := 1
		result := []MonsterSchema{}
		for {
			monsters, err := apiMonsters(page)
			if err != nil {
				return err
			}
			result = append(result, monsters.Data...)
			if page >= *monsters.Pages {
				break
			}
			page++
		}
		output(result)
		return nil
	},
}

func init() {
	databaseCmd.AddCommand(dbMonstersCmd)
}
