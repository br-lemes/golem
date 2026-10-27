package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/spf13/cobra"
)

var reactionCmd = &cobra.Command{
	Args:  cobra.NoArgs,
	Use:   "reaction",
	Short: "Manage event reactions",
}

func init() {
	rootCmd.AddCommand(reactionCmd)
}

func validateReactionCode(code string) error {
	if !slices.Contains(catalog.EventContentCodes(), code) {
		return fmt.Errorf("event content %s not found", code)
	}
	return nil
}
