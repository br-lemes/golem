package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var configFlag string
var outputFlag string
var refreshFlag bool

var rootCmd = &cobra.Command{
	Use:   "golem",
	Short: "A Go CLI to play and automate ArtifactsMMO.",
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
		version := cmd.Root().Version
		if !strings.HasPrefix(version, "v") {
			fmt.Fprintf(cmd.ErrOrStderr(), "%s version %s\n", cmd.Root().Name(), version)
		}

		cfg, err := config.Load(configFlag)
		if err != nil {
			return err
		}
		err = cache.Initialize(cfg.Database)
		if err != nil {
			return fmt.Errorf("initialize cache: %w", err)
		}
		token, prompted := api.Initialize(cfg.API)
		if prompted {
			cfg.API.Token = token
			err = config.Save(configFlag, cfg)
			if err != nil {
				return err
			}
		}

		err = validateFlags()
		if err != nil {
			return err
		}
		if refreshFlag && cmd != refreshCmd {
			return refreshCaches(false)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "Configuration file path")
	rootCmd.PersistentFlags().BoolVarP(&console.Debug, "debug", "d", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&console.Format, "format", "f", "auto", "Output format: auto, json or yaml")
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "", "Output file path (default: stdout)")
	rootCmd.PersistentFlags().BoolVar(&refreshFlag, "refresh", false, "Refresh all caches before running the command")
	rootCmd.PersistentFlags().StringVarP(&console.Style, "style", "s", "monokai", "The style to use for syntax highlighting")
}

func Execute(version string) error {
	rootCmd.Version = version
	err := rootCmd.Execute()
	closer, ok := console.Stdout.(io.Closer)
	if ok && console.Stdout != os.Stdout {
		_ = closer.Close()
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
