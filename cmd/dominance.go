package cmd

import (
	"fmt"
	"sort"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/database"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/surplus"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

type dominanceFlags struct {
	All  bool `flag:"all" desc:"Compare against all equipment in the catalog"`
	By   bool `flag:"by" desc:"List owned equipment that dominates the item"`
	Over bool `flag:"over" desc:"List owned equipment dominated by the item"`
}

var dominanceCmd = &cobra.Command{
	Args:  cobra.ExactArgs(1),
	Use:   "dominance <code>",
	Short: "Analyze equipment dominance",
	Long: `Analyze equipment dominance

Arguments:
  code   The code of the item.`,
	ValidArgsFunction: completion.Item(1).Build(),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]
		flags, err := utils.ReadFlags[dominanceFlags](cmd)
		if err != nil {
			return err
		}
		err = dominanceValidate(code, flags)
		if err != nil {
			return err
		}
		cmd.SilenceUsage = true
		return dominanceRun(code, flags)
	},
}

type dominanceResult struct {
	Code        string                `json:"code"`
	Status      string                `json:"status"`
	Reason      string                `json:"reason,omitempty"`
	DominatedBy []dominanceComparison `json:"dominated_by,omitempty"`
	Dominates   []dominanceComparison `json:"dominates,omitempty"`
}

type dominanceComparison struct {
	Code    string            `json:"code"`
	Level   int               `json:"level"`
	Effects map[string]string `json:"effects,omitempty"`
}

func dominanceValidate(code string, flags dominanceFlags) error {
	if !flags.By && !flags.Over {
		return fmt.Errorf("at least one of --by or --over is required")
	}
	item, exists := database.Items.Get(code)
	if !exists {
		return fmt.Errorf("item not found in catalog: %s", code)
	}
	_, equipment := database.EquipmentTypeToSlots[item.Type]
	if !equipment {
		return fmt.Errorf("item is not equipment: %s", code)
	}
	return nil
}

func dominanceRun(code string, flags dominanceFlags) error {
	item, _ := database.Items.Get(code)
	characters, err := api.AccountsCharacters("")
	if err != nil {
		return err
	}
	bankItems, err := api.MyBankItems()
	if err != nil {
		return err
	}
	owned := ownedEquipment(bankItems, characters)
	candidates := owned
	if flags.All {
		candidates = catalogEquipment()
	}
	result := dominanceResult{Code: code, Status: "equippable"}
	if !canEquipAny(*item, characters) {
		result.Status = "not_equippable"
		result.Reason = "no current character can equip this item"
	}
	for candidateCode, candidate := range candidates {
		if candidateCode == code {
			continue
		}
		if flags.By && surplus.DominatesEquipment(candidate, *item) {
			result.DominatedBy = append(result.DominatedBy, dominanceComparison{
				Code:    candidateCode,
				Level:   candidate.Level,
				Effects: surplus.CompareEffects(*item, candidate),
			})
		}
		if flags.Over && surplus.DominatesEquipment(*item, candidate) {
			result.Dominates = append(result.Dominates, dominanceComparison{
				Code:    candidateCode,
				Level:   candidate.Level,
				Effects: surplus.CompareEffects(*item, candidate),
			})
		}
	}
	sort.Slice(result.DominatedBy, func(i, j int) bool { return result.DominatedBy[i].Code < result.DominatedBy[j].Code })
	sort.Slice(result.Dominates, func(i, j int) bool { return result.Dominates[i].Code < result.Dominates[j].Code })
	return console.Auto(result)
}

func catalogEquipment() map[string]schemas.ItemSchema {
	result := make(map[string]schemas.ItemSchema)
	for _, item := range database.Items.All() {
		_, equipment := database.EquipmentTypeToSlots[item.Type]
		if equipment {
			result[item.Code] = *item
		}
	}
	return result
}

func ownedEquipment(bank []schemas.SimpleItemSchema, characters []schemas.CharacterSchema) map[string]schemas.ItemSchema {
	owned := map[string]schemas.ItemSchema{}
	add := func(code string) {
		item, ok := database.Items.Get(code)
		if ok {
			_, equipment := database.EquipmentTypeToSlots[item.Type]
			if equipment {
				owned[code] = *item
			}
		}
	}
	for _, item := range bank {
		add(item.Code)
	}
	for _, character := range characters {
		if character.Inventory != nil {
			for _, item := range *character.Inventory {
				add(item.Code)
			}
		}
		for _, code := range []string{
			character.AmuletSlot,
			character.Artifact1Slot,
			character.Artifact2Slot,
			character.Artifact3Slot,
			character.BagSlot,
			character.BodyArmorSlot,
			character.BootsSlot,
			character.HelmetSlot,
			character.LegArmorSlot,
			character.Ring1Slot,
			character.Ring2Slot,
			character.ShieldSlot,
			character.Utility1Slot,
			character.Utility2Slot,
			character.WeaponSlot,
		} {
			if code != "" {
				add(code)
			}
		}
	}
	return owned
}

func canEquipAny(item schemas.ItemSchema, characters []schemas.CharacterSchema) bool {
	for _, character := range characters {
		if surplus.CanUse(item, character) {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(dominanceCmd)
	err := utils.RegisterFlags[dominanceFlags](dominanceCmd)
	if err != nil {
		panic(err)
	}
}
