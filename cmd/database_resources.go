package cmd

import (
	"github.com/br-lemes/golem/pkg/console"
	. "github.com/br-lemes/golem/pkg/schemas"
	"github.com/spf13/cobra"
)

var dbResourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		page := 1
		result := []ResourceSchema{}
		for {
			resources, err := apiResources(page)
			if err != nil {
				return err
			}
			result = append(result, resources.Data...)
			if page >= *resources.Pages {
				break
			}
			page++
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbResourcesCmd)
}
