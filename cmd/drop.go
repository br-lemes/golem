package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop",
	RunE: func(cmd *cobra.Command, args []string) error {
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}
		maxLevel := 0
		for _, character := range characters {
			maxLevel = max(maxLevel, character.Level)
		}
		monsters := database.Monsters.Filter(func(m *schemas.MonsterSchema) bool {
			return m.Level <= maxLevel
		})
		drops := map[string]*ItemStock{}
		for _, monster := range monsters {
			for _, drop := range monster.Drops {
				_, exists := drops[drop.Code]
				if !exists {
					drops[drop.Code] = &ItemStock{}
				}
				safety := 1000 / drop.Rate
				drops[drop.Code].Safety = max(drops[drop.Code].Safety, safety)
				drops[drop.Code].Target = max(drops[drop.Code].Target, safety+safety/3)
			}
		}
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		for _, item := range bankItems {
			_, exists := drops[item.Code]
			if !exists {
				continue
			}
			drops[item.Code].Current += item.Quantity
		}
		for _, character := range characters {
			if character.Inventory == nil {
				continue
			}
			for _, slot := range *character.Inventory {
				if slot.Code == "" {
					continue
				}
				_, exists := drops[slot.Code]
				if !exists {
					continue
				}
				drops[slot.Code].Current += slot.Quantity
			}
		}
		for code, drop := range drops {
			if drop.Current >= drop.Safety {
				delete(drops, code)
			} else {
				drops[code].Needed = max(0, drop.Target-drop.Current)
			}
		}
		return console.Auto(drops)
	},
}

func init() {
	rootCmd.AddCommand(dropCmd)
}
