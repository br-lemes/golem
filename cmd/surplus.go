package cmd

import (
	"fmt"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/surplus"
	"github.com/spf13/cobra"
)

var surplusCmd = &cobra.Command{
	Args:  cobra.MaximumNArgs(1),
	Use:   "surplus [code]",
	Short: "List surplus items or show equipment details",
	Long: `List surplus items or show equipment details

Arguments:
  code   The code of the item.`,
	ValidArgsFunction: completion.Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		evaluate, _ := cmd.Flags().GetBool("evaluate")
		bankItems, err := api.MyBankItems()
		if err != nil {
			return err
		}
		characters, err := api.AccountsCharacters("")
		if err != nil {
			return err
		}

		input := surplus.Input{BankItems: bankItems, Characters: characters}

		if len(args) == 1 {
			if evaluate {
				return console.Auto(formatEquipmentExplanation(surplus.Evaluate(input, args[0])))
			}
			return console.Auto(formatEquipmentExplanation(surplus.Explain(input, args[0])))
		}
		if evaluate {
			return fmt.Errorf("--evaluate requires an equipment code")
		}
		items := surplus.Find(input)
		quantities := make(map[string]int, len(items))
		for _, item := range items {
			quantities[item.Item.Code] = item.Surplus
		}
		return console.Auto(quantities)
	},
}

type equipmentExplanation struct {
	Code        string                `json:"code"`
	Status      string                `json:"status"`
	Reason      string                `json:"reason"`
	Level       int                   `json:"level,omitempty"`
	Total       int                   `json:"total,omitempty"`
	Surplus     int                   `json:"surplus,omitempty"`
	DominatedBy []equipmentComparison `json:"dominated_by,omitempty"`
	ComparedTo  []equipmentComparison `json:"compared_to,omitempty"`
}
type equipmentComparison struct {
	Code    string            `json:"code"`
	Level   int               `json:"level"`
	Effects map[string]string `json:"effects,omitempty"`
}

func formatEquipmentExplanation(explanation surplus.Explanation) equipmentExplanation {
	result := equipmentExplanation{
		Code:    explanation.Code,
		Status:  explanation.Status,
		Reason:  explanation.Reason,
		Level:   explanation.Item.Level,
		Total:   explanation.Total,
		Surplus: explanation.Surplus,
	}
	for _, comparison := range explanation.DominatedBy {
		result.DominatedBy = append(result.DominatedBy, equipmentComparison{
			Code:    comparison.Code,
			Level:   comparison.Level,
			Effects: comparison.Effects,
		})
	}
	for _, comparison := range explanation.ComparedTo {
		result.ComparedTo = append(result.ComparedTo, equipmentComparison{
			Code:    comparison.Code,
			Level:   comparison.Level,
			Effects: comparison.Effects,
		})
	}
	return result
}

func init() {
	rootCmd.AddCommand(surplusCmd)
	surplusCmd.Flags().Bool("evaluate", false, "Evaluate an item as if it were owned")
}
