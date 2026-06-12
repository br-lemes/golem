package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
)

var accountsCharactersCmd = &cobra.Command{
	Use:   "accountsCharacters [account]",
	Short: "Get Account Characters",
	Long: `Get Account Characters

Arguments:
  [account]  The name of the account.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 0 {
			return fmt.Errorf("missing required argument: account")
		}
		path = fmt.Sprintf("/accounts/%s/characters", args[0])

		params := make(map[string]string)

		resp, err := apiGet(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(accountsCharactersCmd)
}
