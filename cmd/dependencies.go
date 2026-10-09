package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type dependenciesFlags struct {
	Quantities        bool `flag:"quantities" shorthand:"q" desc:"Include quantities available in the bank"`
	Level             bool `flag:"level" desc:"Include item levels"`
	Paths             bool `flag:"paths" shorthand:"p" desc:"Show only dependency paths"`
	IncludeExactStock bool `flag:"include-exact-stock" desc:"Include items meeting their exact stock quantity"`
}

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
		err = dependenciesValidate(args)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return dependenciesRun(args, flags)
	},
}

func dependenciesValidate(codes []string) error {
	for _, code := range codes {
		_, found := catalog.Items().Get(code)
		if !found {
			return fmt.Errorf("unknown item %q", code)
		}
	}
	return nil
}

func dependenciesRun(codes []string, flags dependenciesFlags) error {
	ignoredStockCodes := map[string]bool{}
	quantities := map[string]int(nil)
	var bankItems []schemas.SimpleItemSchema
	if len(codes) == 0 || flags.Quantities || !flags.IncludeExactStock {
		items, err := api.MyBankItems()
		if err != nil {
			return err
		}
		bankItems = items
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
	if !flags.IncludeExactStock {
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}
		ignoredStockCodes = dependenciesIgnoredStockCodes(cache.ListStocks(), stockItemQuantities(bankItems, characters))
	}
	if len(codes) > 0 {
		codes = slices.DeleteFunc(codes, func(code string) bool {
			return ignoredStockCodes[code]
		})
	}
	return console.Auto(dependencyTree(codes, quantities, !flags.Paths, flags.Level, ignoredStockCodes))
}

func dependenciesIgnoredStockCodes(stocks []models.Stock, quantities map[string]int) map[string]bool {
	result := map[string]bool{}
	for _, stock := range stocks {
		if stock.Enabled && stock.Mode == models.StockModeExact && quantities[stock.Code] >= stock.Quantity {
			result[stock.Code] = true
		}
	}
	return result
}

func dependencyTree(codes []string, quantities map[string]int, expand bool, includeLevels bool, ignoredStockCodes map[string]bool) map[string]any {
	items := catalog.Items().All()
	npcItems := catalog.NpcsItems.All()
	levels := map[string]int(nil)
	if includeLevels {
		levels = make(map[string]int, len(items))
		for _, item := range items {
			levels[item.Code] = item.Level
		}
	}
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
				result[root] = dependencyProducts(root, products, quantities, levels, ignoredStockCodes, map[string]bool{})
			}
			continue
		}
		for _, path := range paths {
			for index, pathCode := range path {
				if ignoredStockCodes[pathCode] {
					path = path[:index]
					break
				}
			}
			addDependencyPath(result, path, quantities, levels)
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

func addDependencyPath(tree map[string]any, path []string, quantities map[string]int, levels map[string]int) {
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
			if levels != nil {
				next["_level"] = levels[code]
			}
			tree[code] = next
		}
		tree = next
	}
}

func dependencyProducts(code string, products map[string][]string, quantities map[string]int, levels map[string]int, ignoredStockCodes map[string]bool, ancestors map[string]bool) map[string]any {
	ancestors[code] = true
	defer delete(ancestors, code)

	result := make(map[string]any)
	if quantities != nil {
		result["_quantity"] = quantities[code]
	}
	if levels != nil {
		result["_level"] = levels[code]
	}
	for _, product := range products[code] {
		if ignoredStockCodes[product] || ancestors[product] {
			continue
		}
		result[product] = dependencyProducts(product, products, quantities, levels, ignoredStockCodes, ancestors)
	}
	return result
}

func init() {
	rootCmd.AddCommand(dependenciesCmd)
	err := utils.RegisterFlags[dependenciesFlags](dependenciesCmd)
	if err != nil {
		panic(err)
	}
}
