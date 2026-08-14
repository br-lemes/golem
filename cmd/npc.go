package cmd

import "github.com/spf13/cobra"

var npcCmd = &cobra.Command{Use: "npc", Short: "NPC"}

func init() {
	rootCmd.AddCommand(npcCmd)
}
