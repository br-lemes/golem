package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var leaderboardCharactersCmd = &cobra.Command{
	Use:   "leaderboardCharacters",
	Short: "Get Characters Leaderboard",
	Args: func(cmd *cobra.Command, args []string) error {
		argCount := len(args)
		switch argCount {
		case 0:
			return nil
		default:
			return fmt.Errorf("invalid number of arguments: %d", argCount)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		argCount := len(args)

		switch argCount {
		case 0:
			path = "/leaderboard/characters"
		}

		params := make(map[string]string)
		localFlags := cmd.LocalFlags()
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if localFlags.Lookup(f.Name) == nil {
				return
			}
			params[f.Name] = f.Value.String()
		})

		resp, err := api.Get(path, params)
		if err != nil {
			return err
		}
		return console.Auto(resp)
	},
}

func init() {
	apiCmd.AddCommand(leaderboardCharactersCmd)
	leaderboardCharactersCmd.Flags().String("name", "", "Character name.")
	leaderboardCharactersCmd.Flags().Int("page", 0, "Page number")
	leaderboardCharactersCmd.Flags().Int("size", 0, "Page size")
	leaderboardCharactersCmd.Flags().String("sort", "", "Sort of character leaderboards.")
}
