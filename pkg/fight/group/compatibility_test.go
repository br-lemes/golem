//go:build groupcompat

package group_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/br-lemes/golem/pkg/catalog"
	"github.com/br-lemes/golem/pkg/fight"
	"github.com/br-lemes/golem/pkg/fight/group"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSoloCompatibility(t *testing.T) {
	cases := []struct {
		name           string
		file           string
		monster        string
		criticalOffset int
	}{
		{"wolf", "group-simulation-solo.json", "wolf", 0},
		{"king slime", "group-simulation-solo.json", "king_slime", 0},
		{"shell", "group-simulation-shell.json", "king_slime", 0},
		{"healing", "group-simulation-healing.json", "wolf", 0},
		{"lifesteal", "group-simulation-lifesteal.json", "wolf", 1},
		{"enchanted mirror", "group-simulation-mirror.json", "king_slime", 0},
		{"burn", "group-simulation-burn.json", "wolf", 0},
		{"frenzy", "group-simulation-frenzy.json", "wolf", 0},
		{"greed", "group-simulation-greed.json", "wolf", -1},
		{"restore", "group-simulation-restore.json", "wolf", -1},
		{"boost", "group-simulation-boost.json", "wolf", -1},
		{"poison", "group-simulation-poison.json", "spider", -1},
		{"antipoison", "group-simulation-antipoison.json", "spider", -1},
		{"monster corrupted", "group-simulation-solo.json", "corrupted_ogre", -1},
		{"monster reconstitution", "group-simulation-solo.json", "lich", -1},
		{"monster void drain", "group-simulation-solo.json", "duskworm", -1},
		{"monster berserker rage", "group-simulation-solo.json", "sea_marauder", -1},
		{"monster protective bubble", "group-simulation-solo.json", "goblin_priestess", -1},
		{"monster healing", "group-simulation-mirror.json", "cultist_emperor", -1},
		{"monster sun shield", "group-simulation-mirror.json", "sonnengott", -1},
		{"monster frenzy", "group-simulation-mirror.json", "rosenblood", 0},
		{"monster greed", "group-simulation-mirror.json", "baby_red_dragon", -1},
		{"monster burn", "group-simulation-solo.json", "imp", -1},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			character := readCharacter(t, testCase.file)
			monster, ok := catalog.Monsters.Get(testCase.monster)
			if !ok {
				t.Fatalf("monster %q not found", testCase.monster)
			}
			sequence := make([]bool, 100)
			if testCase.criticalOffset >= 0 && testCase.criticalOffset < len(sequence) {
				sequence[testCase.criticalOffset] = true
			}
			options := fight.SimulationOptions{
				Iterations:       1,
				RNG:              func() float64 { return 0.99 },
				CriticalSequence: sequence,
			}
			player := fight.FromLoadout(character.Level, fight.CharacterSlots(character), fight.CharacterUtilities(character))
			solo := fight.Simulate(player, *monster, options)
			fighters := []group.Fighter{group.FighterFromFakeCharacter(0, character)}
			combat := group.NewCombat(fighters, group.MonsterFromSchema(*monster), func() float64 { return 0.99 })
			combat.CriticalSequence = sequence
			party := combat.Run(100)
			if solo.Win != party.Win {
				t.Fatalf("win = %v, group win = %v", solo.Win, party.Win)
			}
			if solo.Turns != party.Turns {
				t.Fatalf("turns = %d, group turns = %d", solo.Turns, party.Turns)
			}
			if int(solo.HPRemaining) != int(party.FinalHP[0]) {
				t.Fatalf("final HP = %d, group final HP = %v", solo.HPRemaining, party.FinalHP[0])
			}
		})
	}
}

func readCharacter(t *testing.T, name string) schemas.FakeCharacterSchema {
	t.Helper()
	path := filepath.Join("..", "..", "..", "examples", name)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var characters []schemas.FakeCharacterSchema
	if err = json.Unmarshal(body, &characters); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	if len(characters) != 1 {
		t.Fatalf("%s contains %d characters", path, len(characters))
	}
	return characters[0]
}
