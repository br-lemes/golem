package cmd

import "github.com/spf13/cobra"

var botCmd = &cobra.Command{
	Use:   "bot",
	Short: "Bot commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)
}
