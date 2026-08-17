package cmd

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type findFlags struct {
	From  string `flag:"from" desc:"Use X,Y as the route origin instead of the character's current position"`
	Layer string `flag:"layer" desc:"Use this map layer instead of the character's current layer"`
}

var findCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "find <name> <code>",
	Short: "Find the closest coordinates containing a specific code",
	Long: `Find the closest coordinates containing a specific code

Arguments:
  name   Name of your character.
  code   The target identifier (e.g., chicken, ash_tree, copper_rocks, bank).`,
	ValidArgsFunction: completion.CharacterName(1).Map(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		code := args[1]

		codes := append(database.MapCodes(), database.EventContentCodes()...)
		if !slices.Contains(codes, code) {
			return fmt.Errorf("code '%s' not found", code)
		}
		options, err := utils.ReadFlags[findFlags](cmd)
		if err != nil {
			return err
		}
		if options.Layer != "" && !slices.Contains(database.Enum("MapLayer"), options.Layer) {
			return fmt.Errorf("invalid layer %q: allowed values are %v", options.Layer, database.Enum("MapLayer"))
		}
		character, err := api.Characters(name)
		if err != nil {
			return err
		}
		if options.From != "" {
			character.X, character.Y, err = parseOrigin(options.From)
			if err != nil {
				return err
			}
		}
		if options.Layer != "" {
			character.Layer = schemas.MapLayer(options.Layer)
		}

		tile := database.FindClosest(character, code)
		if tile == nil {
			return fmt.Errorf("no coordinates found for code %s", code)
		}
		return console.Auto(tile)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
	err := utils.RegisterFlags[findFlags](findCmd)
	if err != nil {
		panic(err)
	}
	err = findCmd.RegisterFlagCompletionFunc("layer", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return database.Enum("MapLayer"), cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}

func parseOrigin(value string) (int, int, error) {
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid origin %q: expected X,Y", value)
	}
	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid origin %q: X must be an integer", value)
	}
	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid origin %q: Y must be an integer", value)
	}
	return x, y, nil
}
