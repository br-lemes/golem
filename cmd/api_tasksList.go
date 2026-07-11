package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var tasksListCmd = &cobra.Command{
	Use:   "tasksList [code]",
	Short: "Get All Tasks",
	Long: `Get All Tasks

Arguments:
  code   The code of the task.`,
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 0:
			return nil
		case 1:
			return nil
		default:
			return fmt.Errorf("invalid number of arguments: %d", argCount)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		argCount := len(args)

		switch argCount {
		case 0:
			path = "/tasks/list"
		case 1:
			path = fmt.Sprintf("/tasks/list/%s", args[0])
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
			params[f.Name] = f.Value.String()
		})

		resp, err := api.Get(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(tasksListCmd)
	tasksListCmd.Flags().Int("max_level", 0,
		"Maximum level.")
	tasksListCmd.Flags().Int("min_level", 0,
		"Minimum level.")
	tasksListCmd.Flags().Int("page", 0,
		"Page number")
	tasksListCmd.Flags().Int("size", 0,
		"Page size")
	tasksListCmd.Flags().String("skill", "",
		"Skill of tasks.")
	tasksListCmd.Flags().String("type", "",
		"Type of tasks.")
}
