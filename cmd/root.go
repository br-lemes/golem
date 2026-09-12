package cmd

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/alecthomas/chroma/v2/styles"
	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/config"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/logs"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type persistentFlags struct {
	Config    string   `flag:"config" desc:"Configuration file path"`
	Debug     bool     `flag:"debug" shorthand:"d" desc:"Enable debug mode"`
	Color     bool     `flag:"color" desc:"Force colored output"`
	Exclude   []string `flag:"exclude" desc:"Exclude output paths"`
	ExcludeIf []string `flag:"exclude-if" desc:"Exclude output entries matching conditions"`
	Format    string   `flag:"format" shorthand:"f" default:"auto" desc:"Output format: auto, json or yaml"`
	NoFilters bool     `flag:"no-db-filters" desc:"Ignore output filters from the database"`
	Only      []string `flag:"only" desc:"Keep only output paths"`
	Output    string   `flag:"output" shorthand:"o" desc:"Output file path (default: stdout)"`
	Refresh   bool     `flag:"refresh" desc:"Refresh all caches before running the command"`
	Style     string   `flag:"style" default:"monokai" desc:"The style to use for syntax highlighting"`
}

var rootCmd = &cobra.Command{
	Use:   "golem",
	Short: "A Go CLI to play and automate ArtifactsMMO.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[persistentFlags](cmd)
		if err != nil {
			return err
		}
		console.Debug = flags.Debug
		console.Color = flags.Color
		console.Exclude = flags.Exclude
		console.ExcludeIf = flags.ExcludeIf
		console.Format = flags.Format
		console.Only = flags.Only
		console.Style = flags.Style

		if flags.Output != "" {
			file, err := os.Create(flags.Output)
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
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s version %s\n", cmd.Root().Name(), version)
		}

		cfg, err := config.Load(flags.Config)
		if err != nil {
			return err
		}
		err = cache.Initialize(cfg.Storage)
		if err != nil {
			return fmt.Errorf("initialize cache: %w", err)
		}
		err = logs.Initialize(cfg.Storage)
		if err != nil {
			return fmt.Errorf("initialize logs: %w", err)
		}
		if !flags.NoFilters {
			filters := cache.GetOutputFilters(cmd.CommandPath())
			if !cmd.Flags().Changed("exclude") {
				console.Exclude = filters["exclude"]
			}
			if !cmd.Flags().Changed("exclude-if") {
				console.ExcludeIf = filters["exclude-if"]
			}
			if !cmd.Flags().Changed("only") {
				console.Only = filters["only"]
			}
		}
		if len(console.Exclude) > 0 || len(console.ExcludeIf) > 0 || len(console.Only) > 0 {
			console.Debugf("active output filters: command=%s exclude=%v exclude-if=%v only=%v\n", cmd.CommandPath(), console.Exclude, console.ExcludeIf, console.Only)
		}
		token, prompted := api.Initialize(cfg.API)
		if prompted {
			cfg.API.Token = token
			err = config.Save(flags.Config, cfg)
			if err != nil {
				return err
			}
		}

		if !slices.Contains(console.ValidFormats, console.Format) {
			return fmt.Errorf("invalid format: %s", console.Format)
		}
		if flags.Refresh && cmd != refreshCmd {
			return refreshCaches(false)
		}
		return nil
	},
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

func init() {
	err := utils.RegisterPersistentFlags[persistentFlags](rootCmd)
	if err != nil {
		panic(err)
	}
	err = rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return console.ValidFormats, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
	err = rootCmd.RegisterFlagCompletionFunc("style", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return styles.Names(), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}
