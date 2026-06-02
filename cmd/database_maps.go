package cmd

import (
	. "github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var dbMapsCmd = &cobra.Command{
	Use:   "maps",
	Short: "Maps",
	RunE: func(cmd *cobra.Command, args []string) error {
		page := 1
		result := []MapSchema{}
		for {
			maps, err := apiMaps(page)
			if err != nil {
				return err
			}
			result = append(result, maps.Data...)
			if page >= *maps.Pages {
				break
			}
			page++
		}
		output(result)
		return nil
	},
}

func init() {
	databaseCmd.AddCommand(dbMapsCmd)
}
