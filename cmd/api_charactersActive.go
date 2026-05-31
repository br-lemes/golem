package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var charactersActiveCmd = &cobra.Command{
	Use:   "charactersActive",
	Short: "Get Active Characters",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/characters/active"

		params := make(map[string]string)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			params[f.Name] = f.Value.String()
		})

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return output(resp)
	},
}

func init() {
	apiCmd.AddCommand(charactersActiveCmd)
	charactersActiveCmd.Flags().Int("page", 0,
		"Page number")
	charactersActiveCmd.Flags().Int("size", 0,
		"Page size")
}
