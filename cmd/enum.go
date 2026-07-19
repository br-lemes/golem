package cmd

import (
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/spf13/cobra"
)

var enumCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "enum [name]",
	Short: "Print enum names and details from OpenAPI schema",
	Long: `Print enum names and details from OpenAPI schema

Arguments:
  name   The name of the enum.`,
	ValidArgsFunction: completion.Custom(1, database.GetEnumNames).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			console.Auto(database.GetEnumNames())
			return nil
		}
		console.Auto(database.GetEnum(args[0]))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(enumCmd)
}
