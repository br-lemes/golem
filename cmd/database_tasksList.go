package cmd

import (
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var dbTasksListCmd = &cobra.Command{
	Use:   "tasksList",
	Short: "Tasks List",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := api.TasksList()
		if err != nil {
			return err
		}
		return console.Auto(result)
	},
}

func init() {
	databaseCmd.AddCommand(dbTasksListCmd)
}
