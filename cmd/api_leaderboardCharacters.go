package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var leaderboardCharactersCmd = &cobra.Command{
	Use:   "leaderboardCharacters",
	Short: "Get Characters Leaderboard",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		path = "/leaderboard/characters"

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
	apiCmd.AddCommand(leaderboardCharactersCmd)
	leaderboardCharactersCmd.Flags().String("name", "",
		"Character name.")
	leaderboardCharactersCmd.Flags().Int("page", 0,
		"Page number")
	leaderboardCharactersCmd.Flags().Int("size", 0,
		"Page size")
	leaderboardCharactersCmd.Flags().String("sort", "",
		"Sort of character leaderboards.")
}
