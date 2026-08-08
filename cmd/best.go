package cmd

import "github.com/spf13/cobra"

var bestCmd = &cobra.Command{
	Use:   "best",
	Short: "Find the best equipment for different activities",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(bestCmd)
}
