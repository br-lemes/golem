package cmd

import "github.com/spf13/cobra"

var bestCmd = &cobra.Command{
	Use:   "best",
	Short: "Find the best equipment for different activities",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(bestCmd)
}
