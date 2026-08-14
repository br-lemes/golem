package cmd

import "github.com/spf13/cobra"

var cleanCmd = &cobra.Command{Use: "clean", Short: "Clean the cache"}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
