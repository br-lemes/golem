package routine

import (
	"errors"
	"testing"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFindBankFromStartingPointWithoutAchievements(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
	if results[0].Target.MapId != 334 || results[0].Distance != 5 {
		t.Fatalf("first result = %#v, want map 334 at distance 5", results[0])
	}
	if results[1].Target.MapId != 955 || results[1].Distance != 20 {
		t.Fatalf("second result = %#v, want map 955 at distance 20", results[1])
	}
}

func TestFindRedSlimeAtEqualDistance(t *testing.T) {
	character := schemas.CharacterSchema{X: 2, Y: -1, Layer: "overworld"}
	results := find(deps{}, character, "red_slime", nil, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
	if results[0].Target.MapId != 175 || results[0].Distance != 1 {
		t.Fatalf("first result = %#v, want map 175 at distance 1", results[0])
	}
	if results[1].Target.MapId != 223 || results[1].Distance != 1 {
		t.Fatalf("second result = %#v, want map 223 at distance 1", results[1])
	}
}

func TestFindBankFromStartingPointWithIslandAchievement(t *testing.T) {
	d := deps{
		hasAchievement: func(code string) (bool, error) {
			return code == "secure_the_island", nil
		},
	}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil, eventPointSet{})

	if len(results) != 3 {
		t.Fatalf("find() returned %d results, want 3", len(results))
	}
	island := results[2]
	if island.Target.MapId != 1234 || island.Distance != 21 {
		t.Fatalf("island result = %#v, want map 1234 at distance 21", island)
	}
	if len(island.Requirements) != 2 || island.Requirements[0].Code != "gold" || island.Requirements[0].Value != 1000 {
		t.Fatalf("island requirements = %#v, want gold and achievement", island.Requirements)
	}
}

func TestFindStrangeRocksWithoutActiveEvent(t *testing.T) {
	loadEvents := func() ([]schemas.ActiveEventSchema, error) {
		return nil, nil
	}
	d := deps{eventsActive: loadEvents}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "strange_rocks", nil, eventPointSet{})

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestHasItemUsesCharacterBeforeLoadingBank(t *testing.T) {
	loaded := 0
	character := schemas.CharacterSchema{WeaponSlot: "cultist_cloak"}
	loadBank := func() bool {
		loaded++
		return true
	}
	bank := []schemas.SimpleItemSchema{{Code: "cultist_cloak", Quantity: 1}}

	if !hasItem(character, "cultist_cloak", 1, loadBank, &bank) {
		t.Fatal("hasItem() = false, want true")
	}
	if loaded != 0 {
		t.Fatalf("bank loads = %d, want 0", loaded)
	}
}

func TestHasItemLoadsBankOnlyWhenNeeded(t *testing.T) {
	loaded := 0
	character := schemas.CharacterSchema{}
	loadBank := func() bool {
		loaded++
		return true
	}
	bank := []schemas.SimpleItemSchema{{Code: "cultist_cloak", Quantity: 1}}

	if !hasItem(character, "cultist_cloak", 1, loadBank, &bank) {
		t.Fatal("hasItem() = false, want true")
	}
	if loaded != 1 {
		t.Fatalf("bank loads = %d, want 1", loaded)
	}
}

func TestFindBankWithoutAchievementDependency(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
}

func TestFindBankWithAchievementLoadError(t *testing.T) {
	d := deps{
		hasAchievement: func(string) (bool, error) {
			return false, errors.New("achievements unavailable")
		},
	}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
}

func TestFindStrangeRocksWithActiveEvent(t *testing.T) {
	event, exists := catalog.Events.Get("strange_apparition")
	if !exists || event.Content == nil || len(event.Maps) == 0 {
		t.Fatal("strange apparition event is not in the catalog")
	}
	loadEvents := func() ([]schemas.ActiveEventSchema, error) {
		eventMap := schemas.MapSchema{
			Layer: event.Maps[0].Layer,
			MapId: event.Maps[0].MapId,
			X:     event.Maps[0].X,
			Y:     event.Maps[0].Y,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: event.Content.Code},
			},
		}
		return []schemas.ActiveEventSchema{{Map: eventMap}}, nil
	}
	d := deps{eventsActive: loadEvents}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "strange_rocks", nil, eventPoints(d, "strange_rocks"))

	if len(results) != 1 {
		t.Fatalf("find() returned %d results, want 1", len(results))
	}
	if results[0].Target.MapId != 871 {
		t.Fatalf("event result map id = %d, want 871", results[0].Target.MapId)
	}
}

func TestFindReusesActiveEventsAcrossTeleportStarts(t *testing.T) {
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
	eventLoads := 0
	d := deps{
		eventsActive: func() ([]schemas.ActiveEventSchema, error) {
			eventLoads++
			return []schemas.ActiveEventSchema{{Map: eventMap}}, nil
		},
	}
	character := schemas.CharacterSchema{Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{
		{Code: "forest_bank_potion", Quantity: 1},
		{Code: "recall_potion", Quantity: 1},
	}

	find(d, character, "strange_rocks", potions, eventPoints(d, "strange_rocks"))
	if eventLoads != 1 {
		t.Fatalf("active event loads = %d, want 1", eventLoads)
	}
}

func TestFindStrangeRocksWithoutEventDependency(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "strange_rocks", nil, eventPointSet{})

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestFindStrangeRocksWithEventLoadError(t *testing.T) {
	loadEvents := func() ([]schemas.ActiveEventSchema, error) {
		return nil, errors.New("events unavailable")
	}
	d := deps{eventsActive: loadEvents}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "strange_rocks", nil, eventPointSet{})

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestFindStrangeRocksIgnoresDifferentActiveEventContent(t *testing.T) {
	event, exists := catalog.Events.Get("attacking_the_island")
	if !exists || event.Content == nil || len(event.Maps) == 0 {
		t.Fatal("attacking the island event is not in the catalog")
	}
	loadEvents := func() ([]schemas.ActiveEventSchema, error) {
		eventMap := schemas.MapSchema{
			Layer: event.Maps[0].Layer,
			MapId: event.Maps[0].MapId,
			X:     event.Maps[0].X,
			Y:     event.Maps[0].Y,
			Interactions: schemas.InteractionSchema{
				Content: &schemas.MapContentSchema{Code: event.Content.Code},
			},
		}
		return []schemas.ActiveEventSchema{{Map: eventMap}}, nil
	}
	d := deps{eventsActive: loadEvents}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "strange_rocks", nil, eventPointSet{})

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestFindGoldRocksFromStartingPoint(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}
	results := find(d, character, "gold_rocks", nil, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("Find() returned %d results, want 2", len(results))
	}
	if results[0].Target.MapId != 83 || results[0].Distance != 10 {
		t.Fatalf("first result = %#v, want map 83 at distance 10", results[0])
	}
	if results[1].Target.MapId != 26 || results[1].Distance != 13 {
		t.Fatalf("second result = %#v, want map 26 at distance 13", results[1])
	}
	for index, result := range results {
		if len(result.Requirements) != 0 {
			t.Fatalf("result %d requirements = %#v, want no requirements", index, result.Requirements)
		}
	}
}

func TestFindMithrilRocksFromGoldRocksRegion(t *testing.T) {
	character := schemas.CharacterSchema{X: 5, Y: -4, Layer: "underground"}
	results := find(deps{}, character, "mithril_rocks", nil, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
	if results[0].Target.MapId != 521 || results[0].Distance != 20 {
		t.Fatalf("first result = %#v, want map 521 at distance 20", results[0])
	}
	if results[1].Target.MapId != 467 || results[1].Distance != 22 {
		t.Fatalf("second result = %#v, want map 467 at distance 22", results[1])
	}
	for index, result := range results {
		if len(result.Transitions) != 2 {
			t.Fatalf("result %d transitions = %d, want 2", index, len(result.Transitions))
		}
		if len(result.Requirements) != 0 {
			t.Fatalf("result %d requirements = %#v, want no requirements", index, result.Requirements)
		}
	}
	if results[0].Transitions[0].MapId != 134 || results[0].Transitions[1].MapId != 571 {
		t.Fatalf("transitions = %#v, want maps 134 then 571", results[0].Transitions)
	}
}

func TestFindBankFromBank(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{X: 4, Y: 1, Layer: "overworld"}

	results := find(d, character, "bank", nil, eventPointSet{})

	if len(results) == 0 {
		t.Fatal("find() returned no banks")
	}
	if results[0].Target.MapId != 334 || results[0].Distance != 0 {
		t.Fatalf("first result = %#v, want map 334 at distance 0", results[0])
	}
	if len(results[0].Transitions) != 0 || len(results[0].Requirements) != 0 {
		t.Fatalf("first result path = %#v, want no transitions or costs", results[0])
	}
}

func TestFindInvalidCode(t *testing.T) {
	character := schemas.CharacterSchema{Layer: "overworld"}
	results := find(deps{}, character, "invalid_code", nil, eventPointSet{})

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestFindBankWithForestBankPotion(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{
		{Code: "forest_bank_potion", Quantity: 1},
		{Code: "recall_potion", Quantity: 1},
	}

	results := find(d, character, "bank", potions, eventPointSet{})

	if len(results) != 6 {
		t.Fatalf("find() returned %d results, want 6", len(results))
	}
	var potionResult *Result
	for index := range results {
		if results[index].Potion != nil {
			potionResult = &results[index]
			break
		}
	}
	if potionResult == nil {
		t.Fatal("find() returned no result using a potion")
	}
	if potionResult.Target.MapId != 955 || potionResult.Distance != 5 {
		t.Fatalf("potion result = %#v, want map 955 at distance %d", potionResult, 5)
	}
	if potionResult.Potion == nil || potionResult.Potion.Code != "forest_bank_potion" {
		t.Fatalf("potion result item = %#v, want forest_bank_potion", potionResult.Potion)
	}
}

func TestFindPrefersNearbyWalkingRouteOverTeleportPotion(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{X: 3, Y: 2, Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{
		{Code: "forest_bank_potion", Quantity: 1},
	}

	results := find(d, character, "bank", potions, eventPointSet{})
	if len(results) == 0 {
		t.Fatal("find() returned no bank routes")
	}
	if results[0].Potion != nil || results[0].Distance >= 5 {
		t.Fatalf("closest route = %#v, want short walking route", results[0])
	}
}

func TestFindTeleportRouteLeavesRestrictedAreaThroughTransition(t *testing.T) {
	character := schemas.CharacterSchema{X: 7, Y: 13, Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{{
		Code:     "enchanted_potion",
		Quantity: 1,
	}}
	results := find(deps{}, character, "baby_red_dragon", potions, eventPointSet{})

	for _, result := range results {
		if result.Potion == nil || result.Potion.Code != "enchanted_potion" {
			continue
		}
		if len(result.Requirements) != 0 {
			continue
		}
		for _, transition := range result.Transitions {
			if transition.MapId == 667 {
				return
			}
		}
	}
	t.Fatal("find() returned no free exit from the Enchanted Forest")
}

func TestFindBankIgnoresNonTeleportPotion(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{
		{Code: "small_health_potion", Quantity: 1},
	}

	results := find(d, character, "bank", potions, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
	for _, result := range results {
		if result.Potion != nil {
			t.Fatalf("result potion = %#v, want nil", result.Potion)
		}
	}
}

func TestFindBankIgnoresUnknownPotion(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{{Code: "unknown_potion", Quantity: 1}}

	results := find(d, character, "bank", potions, eventPointSet{})

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
	for _, result := range results {
		if result.Potion != nil {
			t.Fatalf("result potion = %#v, want nil", result.Potion)
		}
	}
}
