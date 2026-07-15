package cmd

import (
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage and interact with game tasks",
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
