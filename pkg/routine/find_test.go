package routine

import (
	"errors"
	"testing"
	"time"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestFindBankFromStartingPointWithoutAchievements(t *testing.T) {
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return nil, nil
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil)

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
	results := find(deps{}, character, "red_slime", nil)

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
	completedAt := time.Now()
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return []schemas.AccountAchievementSchema{{
			Code:        "secure_the_island",
			CompletedAt: &completedAt,
		}}, nil
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil)

	if len(results) != 3 {
		t.Fatalf("find() returned %d results, want 3", len(results))
	}
	island := results[2]
	if island.Target.MapId != 1234 || island.Distance != 21 {
		t.Fatalf("island result = %#v, want map 1234 at distance 21", island)
	}
	if len(island.Costs) != 1 || island.Costs[0].Code != "gold" || island.Costs[0].Value != 1000 {
		t.Fatalf("island costs = %#v, want 1000 gold", island.Costs)
	}
}

func TestFindStrangeRocksWithoutActiveEvent(t *testing.T) {
	loadEvents := func() ([]schemas.ActiveEventSchema, error) {
		return nil, nil
	}
	d := deps{eventsActive: loadEvents}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "strange_rocks", nil)

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestFindBankWithoutAchievementDependency(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil)

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
}

func TestFindBankWithAchievementLoadError(t *testing.T) {
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return nil, errors.New("achievements unavailable")
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "bank", nil)

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

	results := find(d, character, "strange_rocks", nil)

	if len(results) != 1 {
		t.Fatalf("find() returned %d results, want 1", len(results))
	}
	if results[0].Target.MapId != 871 {
		t.Fatalf("event result map id = %d, want 871", results[0].Target.MapId)
	}
}

func TestFindStrangeRocksWithoutEventDependency(t *testing.T) {
	d := deps{}
	character := schemas.CharacterSchema{Layer: "overworld"}

	results := find(d, character, "strange_rocks", nil)

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

	results := find(d, character, "strange_rocks", nil)

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

	results := find(d, character, "strange_rocks", nil)

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestFindGoldRocksFromStartingPoint(t *testing.T) {
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return nil, nil
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{Layer: "overworld"}
	results := find(d, character, "gold_rocks", nil)

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
		if len(result.Costs) != 0 {
			t.Fatalf("result %d costs = %#v, want no costs", index, result.Costs)
		}
	}
}

func TestFindMithrilRocksFromGoldRocksRegion(t *testing.T) {
	character := schemas.CharacterSchema{X: 5, Y: -4, Layer: "underground"}
	results := find(deps{}, character, "mithril_rocks", nil)

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
		if len(result.Costs) != 0 {
			t.Fatalf("result %d costs = %#v, want no costs", index, result.Costs)
		}
	}
	if results[0].Transitions[0].MapId != 134 || results[0].Transitions[1].MapId != 571 {
		t.Fatalf("transitions = %#v, want maps 134 then 571", results[0].Transitions)
	}
}

func TestFindBankFromBank(t *testing.T) {
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return nil, nil
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{X: 4, Y: 1, Layer: "overworld"}

	results := find(d, character, "bank", nil)

	if len(results) == 0 {
		t.Fatal("find() returned no banks")
	}
	if results[0].Target.MapId != 334 || results[0].Distance != 0 {
		t.Fatalf("first result = %#v, want map 334 at distance 0", results[0])
	}
	if len(results[0].Transitions) != 0 || len(results[0].Costs) != 0 {
		t.Fatalf("first result path = %#v, want no transitions or costs", results[0])
	}
}

func TestFindInvalidCode(t *testing.T) {
	character := schemas.CharacterSchema{Layer: "overworld"}
	results := find(deps{}, character, "invalid_code", nil)

	if len(results) != 0 {
		t.Fatalf("find() returned %d results, want 0", len(results))
	}
}

func TestFindBankWithForestBankPotion(t *testing.T) {
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return nil, nil
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{
		{Code: "forest_bank_potion", Quantity: 1},
	}

	results := find(d, character, "bank", potions)

	if len(results) != 4 {
		t.Fatalf("find() returned %d results, want 4", len(results))
	}
	potionResult := results[2]
	if potionResult.Target.MapId != 955 || potionResult.Distance != 0 {
		t.Fatalf("potion result = %#v, want map 955 at distance 0", potionResult)
	}
	if potionResult.Potion == nil || potionResult.Potion.Code != "forest_bank_potion" {
		t.Fatalf("potion result item = %#v, want forest_bank_potion", potionResult.Potion)
	}
}

func TestFindBankIgnoresNonTeleportPotion(t *testing.T) {
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return nil, nil
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{
		{Code: "small_health_potion", Quantity: 1},
	}

	results := find(d, character, "bank", potions)

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
	loadAchievements := func(string) ([]schemas.AccountAchievementSchema, error) {
		return nil, nil
	}
	d := deps{accountsAchievements: loadAchievements}
	character := schemas.CharacterSchema{Layer: "overworld"}
	potions := []schemas.SimpleItemSchema{{Code: "unknown_potion", Quantity: 1}}

	results := find(d, character, "bank", potions)

	if len(results) != 2 {
		t.Fatalf("find() returned %d results, want 2", len(results))
	}
	for _, result := range results {
		if result.Potion != nil {
			t.Fatalf("result potion = %#v, want nil", result.Potion)
		}
	}
}
