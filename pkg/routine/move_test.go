package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMoveUsesFirstFreePath(t *testing.T) {
	character := schemas.CharacterSchema{Name: "hero", Layer: "overworld"}
	var movedTo []struct{ x, y int }
	d := deps{
		myActionMove: func(_ string, x, y int) (schemas.CharacterMovementDataSchema, error) {
			movedTo = append(movedTo, struct{ x, y int }{x, y})
			character.X = x
			character.Y = y
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
	}

	got, err := move(d, character, "bank", MoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(movedTo) != 1 {
		t.Fatalf("move calls = %d, want 1", len(movedTo))
	}
	if got.X != movedTo[0].x || got.Y != movedTo[0].y {
		t.Fatalf("returned character = %#v, want position %#v", got, movedTo[0])
	}
}

func TestMoveToCurrentStaticDestinationDoesNothing(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     4,
		Y:     1,
		Layer: "overworld",
	}
	d := deps{
		eventsActive: func() ([]schemas.ActiveEventSchema, error) {
			t.Fatal("events loaded for a static destination")
			return nil, nil
		},
		hasAchievement: func(string) (bool, error) {
			t.Fatal("achievements loaded while already at bank")
			return false, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			t.Fatal("movement action called while already at bank")
			return schemas.CharacterMovementDataSchema{}, nil
		},
	}

	got, err := move(d, character, "bank", MoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got != character {
		t.Fatalf("move() = %#v, want unchanged character %#v", got, character)
	}
}

func TestMoveToCurrentEventDestinationChecksActiveEvents(t *testing.T) {
	event, exists := catalog.Events.Get("strange_apparition")
	if !exists || event.Content == nil || len(event.Maps) == 0 {
		t.Fatal("strange apparition event is not in the catalog")
	}
	eventMap := schemas.MapSchema{
		Layer: event.Maps[0].Layer,
		MapId: event.Maps[0].MapId,
		X:     event.Maps[0].X,
		Y:     event.Maps[0].Y,
		Interactions: schemas.InteractionSchema{
			Content: &schemas.MapContentSchema{Code: event.Content.Code},
		},
	}
	character := schemas.CharacterSchema{
		X:     eventMap.X,
		Y:     eventMap.Y,
		Layer: eventMap.Layer,
	}
	eventLoads := 0
	d := deps{
		eventsActive: func() ([]schemas.ActiveEventSchema, error) {
			eventLoads++
			return []schemas.ActiveEventSchema{{Map: eventMap}}, nil
		},
		hasAchievement: func(string) (bool, error) {
			t.Fatal("achievements loaded while already at event target")
			return false, nil
		},
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			t.Fatal("movement action called while already at event target")
			return schemas.CharacterMovementDataSchema{}, nil
		},
	}

	got, err := move(d, character, "strange_rocks", MoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got != character {
		t.Fatalf("move() = %#v, want unchanged character %#v", got, character)
	}
	if eventLoads != 1 {
		t.Fatalf("active event loads = %d, want 1", eventLoads)
	}
}

func TestMakeMoveSkipsAPIWhenAlreadyAtTarget(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     4,
		Y:     1,
		Layer: "overworld",
	}
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			t.Fatal("movement action called at target coordinates")
			return schemas.CharacterMovementDataSchema{}, nil
		},
	}
	got, err := makeMove(d, character, schemas.MapSchema{
		X:     4,
		Y:     1,
		Layer: "overworld",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != character {
		t.Fatalf("makeMove() = %#v, want unchanged character %#v", got, character)
	}
}

func TestMakeMoveDoesNotTreatSameCoordinatesOnDifferentLayerAsTarget(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     4,
		Y:     1,
		Layer: "overworld",
	}
	moves := 0
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			moves++
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
	}
	_, err := makeMove(d, character, schemas.MapSchema{
		X:     4,
		Y:     1,
		Layer: "interior",
	})
	if err != nil {
		t.Fatal(err)
	}
	if moves != 1 {
		t.Fatalf("movement calls = %d, want 1 for a different layer", moves)
	}
}

func TestMoveReturnsErrorWithoutFreePath(t *testing.T) {
	character := schemas.CharacterSchema{Layer: "overworld"}

	_, err := move(deps{}, character, "invalid_code", MoveOptions{})
	if err == nil {
		t.Fatal("move() returned nil error, want error")
	}
}

func TestMoveExecutesTransitions(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     5,
		Y:     -4,
		Layer: "underground",
	}
	moves := 0
	transitions := 0
	d := deps{
		myActionMove: func(_ string, x, y int) (schemas.CharacterMovementDataSchema, error) {
			moves++
			character.X = x
			character.Y = y
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(_ string) (schemas.CharacterTransitionDataSchema, error) {
			transitions++
			return schemas.CharacterTransitionDataSchema{Character: character}, nil
		},
	}

	_, err := move(d, character, "mithril_rocks", MoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if moves != 3 {
		t.Fatalf("move calls = %d, want 3", moves)
	}
	if transitions != 2 {
		t.Fatalf("transition calls = %d, want 2", transitions)
	}
}

func TestMoveUsesTeleportPotionForDistantTarget(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{{
		Code:     "forest_bank_potion",
		Quantity: 1,
	}}
	character := schemas.CharacterSchema{
		Name:              "hero",
		X:                 -5,
		Y:                 -5,
		Layer:             "overworld",
		Inventory:         &inventory,
		InventoryMaxItems: 10,
	}
	used := false
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return nil, nil
		},
		myActionUse: func(_ string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
			used = true
			if item.Code != "forest_bank_potion" || item.Quantity != 1 {
				t.Fatalf("used item = %#v, want one forest bank potion", item)
			}
			character.X = 7
			character.Y = 13
			inventory[0].Quantity = 0
			return schemas.UseItemSchema{Character: character}, nil
		},
	}

	got, err := move(d, character, "bank", MoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !used {
		t.Fatal("teleport potion was not used")
	}
	if got.X != 7 || got.Y != 13 {
		t.Fatalf("position after teleport = (%d, %d), want (7, 13)", got.X, got.Y)
	}
}

func TestAvailableTeleportPotionsCanBeDisabled(t *testing.T) {
	inventory := []schemas.InventorySlotSchema{{
		Code:     "forest_bank_potion",
		Quantity: 1,
	}}
	character := schemas.CharacterSchema{Inventory: &inventory}

	potions, _, _, err := availableTeleportPotions(deps{}, character, MoveOptions{
		NoTeleport: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(potions) != 0 {
		t.Fatalf("availableTeleportPotions() = %#v, want none", potions)
	}
}

func TestMoveNoTeleportSkipsReserveRestock(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:  "hero",
		X:     7,
		Y:     13,
		Layer: "overworld",
	}
	withdrawCalled := false
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "recall_potion", Quantity: 1},
			}, nil
		},
		myActionBankWithdrawItem: func(string, []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			withdrawCalled = true
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
	}

	_, err := move(d, character, "bank", MoveOptions{NoTeleport: true})
	if err != nil {
		t.Fatal(err)
	}
	if withdrawCalled {
		t.Fatal("reserve was refilled with --no-teleport")
	}
}

func TestMoveWithdrawsTeleportPotionFromCurrentBank(t *testing.T) {
	character := schemas.CharacterSchema{
		Name:              "hero",
		X:                 7,
		Y:                 13,
		Layer:             "overworld",
		InventoryMaxItems: 20,
	}
	withdrawn := false
	used := false
	d := deps{
		myBankItems: func() ([]schemas.SimpleItemSchema, error) {
			return []schemas.SimpleItemSchema{
				{Code: "enchanted_potion", Quantity: 1},
			}, nil
		},
		myActionBankWithdrawItem: func(_ string, items []schemas.SimpleItemSchema) (schemas.BankItemTransactionSchema, error) {
			withdrawn = true
			if len(items) != 1 || items[0].Code != "enchanted_potion" || items[0].Quantity != 1 {
				t.Fatalf("withdrawn items = %#v, want one enchanted potion", items)
			}
			inventory := []schemas.InventorySlotSchema{
				{Code: "enchanted_potion", Quantity: 1},
			}
			character.Inventory = &inventory
			return schemas.BankItemTransactionSchema{Character: character}, nil
		},
		myActionUse: func(_ string, item schemas.SimpleItemSchema) (schemas.UseItemSchema, error) {
			used = true
			if item.Code != "enchanted_potion" {
				t.Fatalf("used item = %#v, want enchanted potion", item)
			}
			character.X = -5
			character.Y = 9
			character.Inventory = nil
			return schemas.UseItemSchema{Character: character}, nil
		},
	}

	got, err := move(d, character, "enchanted_mushroom", MoveOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !withdrawn || !used {
		t.Fatalf("withdrawn = %t, used = %t, want both", withdrawn, used)
	}
	if got.X != -5 || got.Y != 9 {
		t.Fatalf("position = (%d, %d), want (-5, 9)", got.X, got.Y)
	}
}

func TestMoveReturnsTransitionError(t *testing.T) {
	character := schemas.CharacterSchema{X: 5, Y: -4, Layer: "underground"}
	wantErr := errors.New("transition failed")
	d := deps{
		myActionMove: func(_ string, x, y int) (schemas.CharacterMovementDataSchema, error) {
			character.X = x
			character.Y = y
			return schemas.CharacterMovementDataSchema{Character: character}, nil
		},
		myActionTransition: func(string) (schemas.CharacterTransitionDataSchema, error) {
			return schemas.CharacterTransitionDataSchema{}, wantErr
		},
	}

	_, err := move(d, character, "mithril_rocks", MoveOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("move() error = %v, want %v", err, wantErr)
	}
}

func TestMoveReturnsMoveErrorAtTransition(t *testing.T) {
	character := schemas.CharacterSchema{X: 5, Y: -4, Layer: "underground"}
	wantErr := errors.New("move failed")
	d := deps{
		myActionMove: func(string, int, int) (schemas.CharacterMovementDataSchema, error) {
			return schemas.CharacterMovementDataSchema{}, wantErr
		},
	}

	_, err := move(d, character, "mithril_rocks", MoveOptions{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("move() error = %v, want %v", err, wantErr)
	}
}

func TestCanUseRequirementsGold(t *testing.T) {
	requirements := []schemas.ConditionSchema{
		{Code: "gold", Operator: "cost", Value: 1000},
	}
	if canUseRequirements(requirements, false) {
		t.Fatal("canUseRequirements() = true without gold permission")
	}
	if !canUseRequirements(requirements, true) {
		t.Fatal("canUseRequirements() = false with gold permission")
	}
}

func TestCanUseRequirementsItemCost(t *testing.T) {
	requirements := []schemas.ConditionSchema{
		{Code: "priestess_hideout_key", Operator: "cost", Value: 1},
	}
	if !canUseRequirements(requirements, false) {
		t.Fatal("item cost should be usable without gold permission")
	}
}

func TestCanUseRequirementsAchievement(t *testing.T) {
	requirements := []schemas.ConditionSchema{
		{Code: "secure_the_island", Operator: "achievement_unlocked", Value: 1},
	}
	if !canUseRequirements(requirements, false) {
		t.Fatal("canUseRequirements() = false for unlocked achievement")
	}
}

func TestRequiredGoldSumsTransitionCosts(t *testing.T) {
	requirements := []schemas.ConditionSchema{
		{Code: "gold", Operator: "cost", Value: 1000},
		{Code: "gold", Operator: "cost", Value: 5000},
		{Code: "key", Operator: "has_item", Value: 1},
	}
	got := requiredGold(requirements)
	if got != 6000 {
		t.Fatalf("requiredGold() = %d, want 6000", got)
	}
}

func TestMissingItemsIncludesItemCosts(t *testing.T) {
	character := schemas.CharacterSchema{
		Inventory: &[]schemas.InventorySlotSchema{{
			Code:     "priestess_hideout_key",
			Quantity: 1,
		}},
	}
	requirements := []schemas.ConditionSchema{
		{Code: "priestess_hideout_key", Operator: "cost", Value: 2},
		{Code: "gold", Operator: "cost", Value: 1000},
	}
	got := missingItems(character, requirements)
	valid := len(got) == 1 && got[0].Code == "priestess_hideout_key" && got[0].Quantity == 1
	if !valid {
		t.Fatalf("missingItems() = %#v, want one missing key", got)
	}
}

func TestMissingItemCostDoesNotCountEquippedItems(t *testing.T) {
	character := schemas.CharacterSchema{WeaponSlot: "priestess_hideout_key"}
	requirements := []schemas.ConditionSchema{
		{Code: "priestess_hideout_key", Operator: "cost", Value: 1},
	}
	got := missingItems(character, requirements)
	if len(got) != 1 || got[0].Code != "priestess_hideout_key" || got[0].Quantity != 1 {
		t.Fatalf("missingItems() = %#v, want one key from inventory or bank", got)
	}
}
