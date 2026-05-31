package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var resourcesCmd = &cobra.Command{
	Use:   "resources [code]",
	Short: "Get All Resources",
	Long: `Get All Resources

Arguments:
  [code]  The code of the resource.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			path = "/resources"
		} else if len(args) == 1 {
			path = fmt.Sprintf("/resources/%s", args[0])
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
	apiCmd.AddCommand(resourcesCmd)
	resourcesCmd.Flags().String("drop", "",
		"Item code of the drop.")
	resourcesCmd.Flags().Int("max_level", 0,
		"Maximum level.")
	resourcesCmd.Flags().Int("min_level", 0,
		"Minimum level.")
	resourcesCmd.Flags().Int("page", 0,
		"Page number")
	resourcesCmd.Flags().Int("size", 0,
		"Page size")
	resourcesCmd.Flags().String("skill", "",
		"Skill of resources.")
}
