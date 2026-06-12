package cmd

import (
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var myCharactersCmd = &cobra.Command{
	Use:   "myCharacters",
	Short: "Get My Characters",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/my/characters"

		params := make(map[string]string)

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(myCharactersCmd)
}
