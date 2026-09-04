package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

var filterKinds = []string{"exclude", "exclude-if", "only"}

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Manage persistent output filters",
}

func validateFilterKind(kind string) error {
	if !slices.Contains(filterKinds, kind) {
		return fmt.Errorf("invalid filter kind %q: allowed values are %v", kind, filterKinds)
	}
	return nil
}

func validateFilterArgs(args []string, count int) error {
	if len(args) != count {
		return fmt.Errorf("requires %d arguments", count)
	}
	if args[0] == "" {
		return fmt.Errorf("command cannot be empty")
	}
	err := validateFilterKind(args[1])
	if err != nil {
		return err
	}
	if args[2] == "" {
		return fmt.Errorf("pattern cannot be empty")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(filterCmd)
}
