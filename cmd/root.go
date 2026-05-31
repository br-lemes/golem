package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/br-lemes/golem/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	debugFlag  bool
	formatFlag string
	outputFlag string
	styleFlag  string
	userToken  string
	writer     io.Writer
)

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
			writer = file
		} else {
			writer = cmd.OutOrStdout()
		}

		err := validateFlags()
		if err != nil {
			return err
		}
		err = loadConfig()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false,
		"Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "auto",
		"Formato da saída: auto, json ou yaml")
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "",
		"Output file path (default: stdout)")
	rootCmd.PersistentFlags().StringVarP(&styleFlag, "style", "s", "monokai",
		"The style to use for syntax highlighting")
}

func Execute() error {
	err := rootCmd.Execute()
	closer, ok := writer.(io.Closer)
	if ok && writer != os.Stdout {
		closer.Close()
	}
	return err
}

func validateFlags() error {
	switch formatFlag {
	case "auto", "json", "yaml":
	default:
		return fmt.Errorf("invalid format: %s", formatFlag)
	}
	return nil
}

func loadConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}

	configName := "user"
	configType := "json"
	configFile := filepath.Join(
		configDir,
		fmt.Sprintf("%s.%s", configName, configType),
	)
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(configDir, 0755)
		err := viper.WriteConfigAs(configFile)
		if err != nil {
			return err
		}
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	userToken = viper.GetString("token")
	if userToken == "" {
		fmt.Fprintf(os.Stdout, "Enter your token: ")
		_, err := fmt.Scan(&userToken)
		if err != nil {
			return err
		}
		viper.Set("token", userToken)
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
	}

	return nil
}

func getConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "golem"), nil
}
