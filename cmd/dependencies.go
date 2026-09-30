package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type dependenciesFlags struct {
	Quantities bool   `flag:"quantities" shorthand:"q" desc:"Include quantities available in the bank"`
	Mode       string `flag:"mode" shorthand:"m" default:"auto" desc:"Output mode: auto, expand, or paths"`
}

var dependencyModes = []string{"auto", "expand", "paths"}

var dependenciesCmd = &cobra.Command{
	Args:  cobra.ArbitraryArgs,
	Use:   "dependencies [code]...",
	Short: "Show what items and currencies can be used for",
	Long: `Show what items and currencies can be used for

Arguments:
  code   The code of the item.`,
	ValidArgsFunction: completion.Item(0).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags, err := utils.ReadFlags[dependenciesFlags](cmd)
		if err != nil {
			return err
		}
		err = dependenciesValidate(flags)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return dependenciesRun(args, flags)
	},
}

func dependenciesValidate(flags dependenciesFlags) error {
	if !slices.Contains(dependencyModes, flags.Mode) {
		return fmt.Errorf("invalid mode %q: allowed values are %v", flags.Mode, dependencyModes)
	}
	return nil
}

func dependenciesRun(codes []string, flags dependenciesFlags) error {
	expand := dependenciesExpand(flags.Mode, len(codes) > 0)
	quantities := map[string]int(nil)
	if len(codes) == 0 || flags.Quantities {
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		if flags.Quantities {
			quantities = make(map[string]int, len(bankItems))
			for _, item := range bankItems {
				quantities[item.Code] += item.Quantity
			}
		}
		if len(codes) == 0 {
			for _, item := range bankItems {
				codes = append(codes, item.Code)
			}
		}
	}
	return console.Auto(dependencyTree(codes, quantities, expand))
}

func dependenciesExpand(mode string, hasCodes bool) bool {
	switch mode {
	case "auto":
		return !hasCodes
	case "expand":
		return true
	default:
		return false
	}
}

func dependencyTree(codes []string, quantities map[string]int, expand bool) map[string]any {
	items := catalog.Items().All()
	npcItems := catalog.NpcsItems.All()
	products := make(map[string][]string)
	ingredients := make(map[string][]string)
	for _, item := range items {
		if item.Craft == nil || item.Craft.Items == nil {
			continue
		}
		for _, ingredient := range *item.Craft.Items {
			products[ingredient.Code] = append(products[ingredient.Code], item.Code)
			ingredients[item.Code] = append(ingredients[item.Code], ingredient.Code)
		}
	}
	for _, item := range npcItems {
		if item.BuyPrice == nil || *item.BuyPrice <= 0 {
			continue
		}
		_, currencyIsItem := catalog.Items().Get(item.Currency)
		if !currencyIsItem {
			continue
		}
		products[item.Currency] = append(products[item.Currency], item.Code)
		ingredients[item.Code] = append(ingredients[item.Code], item.Currency)
	}
	for code := range products {
		slices.Sort(products[code])
		products[code] = slices.Compact(products[code])
	}
	for code := range ingredients {
		slices.Sort(ingredients[code])
		ingredients[code] = slices.Compact(ingredients[code])
	}

	result := make(map[string]any, len(codes))
	for _, code := range codes {
		paths := dependencyPaths(code, ingredients, map[string]bool{})
		if expand {
			for _, path := range paths {
				root := path[0]
				result[root] = dependencyProducts(root, products, quantities, map[string]bool{})
			}
			continue
		}
		for _, path := range paths {
			addDependencyPath(result, path, quantities)
		}
	}
	return result
}

func dependencyPaths(code string, ingredients map[string][]string, ancestors map[string]bool) [][]string {
	if ancestors[code] {
		return [][]string{{code}}
	}
	parents := ingredients[code]
	if len(parents) == 0 {
		return [][]string{{code}}
	}
	ancestors[code] = true
	defer delete(ancestors, code)

	var paths [][]string
	for _, parent := range parents {
		for _, path := range dependencyPaths(parent, ingredients, ancestors) {
			paths = append(paths, append(path, code))
		}
	}
	return paths
}

func addDependencyPath(tree map[string]any, path []string, quantities map[string]int) {
	seen := make(map[string]bool, len(path))
	for _, code := range path {
		if seen[code] {
			return
		}
		seen[code] = true

		next, ok := tree[code].(map[string]any)
		if !ok {
			next = make(map[string]any)
			if quantities != nil {
				next["_quantity"] = quantities[code]
			}
			tree[code] = next
		}
		tree = next
	}
}

func dependencyProducts(code string, products map[string][]string, quantities map[string]int, ancestors map[string]bool) map[string]any {
	ancestors[code] = true
	defer delete(ancestors, code)

	result := make(map[string]any)
	if quantities != nil {
		result["_quantity"] = quantities[code]
	}
	for _, product := range products[code] {
		if ancestors[product] {
			continue
		}
		result[product] = dependencyProducts(product, products, quantities, ancestors)
	}
	return result
}

func init() {
	rootCmd.AddCommand(dependenciesCmd)
	err := utils.RegisterFlags[dependenciesFlags](dependenciesCmd)
	if err != nil {
		panic(err)
	}
	err = dependenciesCmd.RegisterFlagCompletionFunc("mode", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return dependencyModes, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}
