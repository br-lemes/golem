package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var tasksListCmd = &cobra.Command{
	Use:   "tasksList [code]",
	Short: "Get All Tasks",
	Long: `Get All Tasks

Arguments:
  [code]  The code of the task.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			path = "/tasks/list"
		} else if len(args) == 1 {
			path = fmt.Sprintf("/tasks/list/%s", args[0])
		}

		params := make(map[string]string)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			params[f.Name] = f.Value.String()
		})

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return output(resp)
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
