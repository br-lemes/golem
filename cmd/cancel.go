package cmd

import (
	"fmt"
	"slices"

	"github.com/br-lemes/golem/pkg/api"
	"github.com/br-lemes/golem/pkg/completion"
	"github.com/br-lemes/golem/pkg/routine"
	"github.com/br-lemes/golem/pkg/schemas"
	"github.com/br-lemes/golem/pkg/utils"
	"github.com/spf13/cobra"
)

var cancelData struct {
	character schemas.CharacterSchema
	order     schemas.GEOrderSchema
}

var cancelCmd = &cobra.Command{
	Args:  cobra.ExactArgs(2),
	Use:   "cancel <name> <id>",
	Short: "Cancel an existing GE order",
	Long: `Cancel an existing GE order

Arguments:
  name   Name of your character.
  id     The id of the order you want to cancel.`,
	ValidArgsFunction: completion.CharacterName(1).Build(),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		id := args[1]
		validCharacters := utils.GetCharacters()
		if !slices.Contains(validCharacters, name) {
			return fmt.Errorf("invalid character %q: allowed values are %v", name, validCharacters)
		}
		if id == "" {
			return fmt.Errorf("id must not be empty")
		}
		orders, err := api.MyGrandexchangeOrders("", "")
		if err != nil {
			return err
		}
		var found *schemas.GEOrderSchema
		for _, order := range orders {
			if order.Id == id {
				found = &order
				break
			}
		}
		if found == nil {
			return fmt.Errorf("order %q not found", id)
		}
		cancelData.order = *found
		cancelData.character, err = api.Characters(name)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		name := args[0]
		id := args[1]

		var err error
		cancelData.character, err = routine.Move(cancelData.character, "grand_exchange")
		if err != nil {
			return err
		}

		cancel := schemas.GECancelOrderSchema{Id: id}
		cancelOrderData, err := api.MyActionGrandexchangeCancel(name, cancel)
		if err != nil {
			return err
		}
		cancelData.character = cancelOrderData.Character
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)
}
