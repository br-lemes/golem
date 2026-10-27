package cache

import (
	"testing"
)

func TestReactions(t *testing.T) {
	initializeTestCache(t)
	setReaction := func(name, code string, priority int) error {
		enabled := true
		return SetReaction(name, code, &priority, &enabled)
	}

	err := setReaction("br5", "bandit_lizard", 2)
	if err != nil {
		t.Fatal(err)
	}
	err = setReaction("br5", "strange_rocks", 1)
	if err != nil {
		t.Fatal(err)
	}
	err = setReaction("br6", "bandit_lizard", 0)
	if err != nil {
		t.Fatal(err)
	}

	reactions := ListReactions([]string{"br5"}, []string{"bandit_lizard"})
	if len(reactions) != 1 {
		t.Fatalf("reactions = %#v, want 1 reaction", reactions)
	}
	if reactions[0].Priority != 2 {
		t.Fatalf("reactions are not ordered by priority: %#v", reactions)
	}
	if len(ListReactions([]string{"br5"}, nil)) != 2 {
		t.Fatal("character filter returned the wrong reactions")
	}
	allReactions := ListReactions(nil, nil)
	if len(allReactions) != 3 || allReactions[0].Name != "br5" || allReactions[2].Name != "br6" {
		t.Fatalf("reactions are not ordered by name: %#v", allReactions)
	}
	if len(ListReactions([]string{"br6"}, []string{"bandit_lizard"})) != 1 {
		t.Fatal("character and code filters returned the wrong reactions")
	}
	if len(ListReactions(nil, []string{"bandit_lizard"})) != 2 {
		t.Fatal("code filter returned the wrong reactions")
	}
	if len(ListReactions([]string{"br5"}, []string{"missing"})) != 0 {
		t.Fatal("found reactions for missing code")
	}
	err = RemoveReactions("bandit_lizard", []string{"br5"})
	if err != nil {
		t.Fatal(err)
	}
	if len(ListReactions([]string{"br5"}, []string{"bandit_lizard"})) != 0 {
		t.Fatal("reaction was not removed")
	}
	err = setReaction("br5", "full_moon_vampire", 1)
	if err == nil {
		t.Fatal("duplicate character priority was accepted")
	}

	err = SetReaction("br5", "strange_rocks", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	reaction, found := GetReaction("br5", "strange_rocks")
	if !found || reaction.Priority != 1 || !reaction.Enabled {
		t.Fatalf("set without options changed the reaction: %#v", reaction)
	}
	err = SetReactionsEnabled("bandit_lizard", nil, true)
	if err != nil {
		t.Fatal(err)
	}
	reaction, found = GetReaction("br6", "bandit_lizard")
	if !found || !reaction.Enabled {
		t.Fatalf("enable all did not enable the reaction: %#v", reaction)
	}
	err = SetReactionsEnabled("bandit_lizard", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	reaction, found = GetReaction("br6", "bandit_lizard")
	if !found || reaction.Enabled {
		t.Fatalf("disable all did not disable the reaction: %#v", reaction)
	}
	err = RemoveReactions("bandit_lizard", nil)
	if err != nil {
		t.Fatal(err)
	}
	_, found = GetReaction("br6", "bandit_lizard")
	if found {
		t.Fatal("remove all did not remove the reaction")
	}
	err = SetReaction("br5", "full_moon_vampire", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	reaction, found = GetReaction("br5", "full_moon_vampire")
	if !found || reaction.Priority != 2 || !reaction.Enabled {
		t.Fatalf("set did not assign the next priority: %#v", reaction)
	}

	CleanReactions()
	if len(ListReactions(nil, nil)) != 0 {
		t.Fatal("event reactions were not cleaned")
	}
}
