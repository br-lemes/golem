package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var routesCmd = &cobra.Command{
	Use:   "routes [route]",
	Short: "List or inspect API routes from the game's OpenAPI spec",
	Long: `List or inspect API routes from the game's OpenAPI spec

Arguments:
  route   Path of a specific route to inspect.`,
	Args: cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return utils.GetRoutesCompletion(),
				cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			res, err := utils.GetRoute(args[0])
			if err != nil {
				return fmt.Errorf("failed to process route: %w", err)
			}
			return console.Auto(res)
		}

		routes, err := utils.GetRoutes()
		if err != nil {
			return fmt.Errorf("failed to extract routes: %w", err)
		}

		return console.Auto(routes)
	},
}

func init() {
	rootCmd.AddCommand(routesCmd)
}
