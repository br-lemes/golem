package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/br-lemes/golem/internal/version"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var outputFlag string

var rootCmd = &cobra.Command{
	Use:     "golem",
	Short:   "A Go CLI to play and automate ArtifactsMMO.",
	Version: version.GetVersion(),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if outputFlag != "" {
			file, err := os.Create(outputFlag)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			console.Stdout = file
			cmd.SetOut(file)
		} else {
			console.Stdout = cmd.OutOrStdout()
		}
		console.Stderr = cmd.ErrOrStderr()
		console.Stdin = cmd.InOrStdin()

		err := validateFlags()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&console.Debug, "debug", "d", false,
		"Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&console.Format, "format", "f", "auto",
		"Formato da saída: auto, json ou yaml")
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "",
		"Output file path (default: stdout)")
	rootCmd.PersistentFlags().StringVarP(&console.Style, "style", "s", "monokai",
		"The style to use for syntax highlighting")
}

func Execute() error {
	err := rootCmd.Execute()
	closer, ok := console.Stdout.(io.Closer)
	if ok && console.Stdout != os.Stdout {
		closer.Close()
	}
	return err
}

func validateFlags() error {
	switch console.Format {
	case "auto", "json", "yaml":
	default:
		return fmt.Errorf("invalid format: %s", console.Format)
	}
	return nil
}
