package cmd

import "github.com/spf13/cobra"

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
