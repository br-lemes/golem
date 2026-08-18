package utils

import (
	"reflect"
	"testing"
)

func TestNormalizeBestPriorities(t *testing.T) {
	got, err := NormalizeBestPriorities([]string{
		"wisdom",
		"mining",
		"prospecting",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"mining", "wisdom", "prospecting", "inventory_space"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %#v, want %#v", got, want)
	}
}

func TestNormalizeBestPrioritiesPreservesExplicitInventoryPosition(t *testing.T) {
	got, err := NormalizeBestPriorities([]string{
		"wisdom",
		"inventory_space",
		"prospecting",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"wisdom", "inventory_space", "prospecting"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %#v, want %#v", got, want)
	}
}

func TestNormalizeBestPrioritiesRejectsInvalidEffect(t *testing.T) {
	_, err := NormalizeBestPriorities([]string{"wisdom", "unknown"})
	if err == nil {
		t.Fatal("expected invalid effect error")
	}
}

func TestNormalizeBestPrioritiesRejectsDuplicateEffect(t *testing.T) {
	_, err := NormalizeBestPriorities([]string{"wisdom", "wisdom"})
	if err == nil {
		t.Fatal("expected duplicate effect error")
	}
}

func TestNormalizeBestPrioritiesRejectsMultipleGatheringSkills(t *testing.T) {
	_, err := NormalizeBestPriorities([]string{"mining", "woodcutting"})
	if err == nil {
		t.Fatal("expected multiple gathering skills error")
	}
}
