package cmd

import (
	_ "embed"
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate <command>",
	Short: "Generate Go code for specific commands",
	Long: `Generate Go source code for specific game API commands

Arguments:
  command   Name of the target command to generate code for.`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return utils.GetCommandsCompletion(),
				cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		targetCmd := args[0]

		data, err := utils.BuildTemplateData(targetCmd)
		if err != nil {
			return err
		}

		formattedCode, err := utils.RenderTemplate(data)
		if err != nil {
			return err
		}

		_, err = console.Stdout.Write(formattedCode)
		if err != nil {
			return fmt.Errorf("failed to write generated code: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
