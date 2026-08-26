package best

import (
	"reflect"
	"testing"
)

func TestNormalizePriorities(t *testing.T) {
	got, err := NormalizePriorities([]string{"wisdom", "mining", "prospecting"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"mining", "wisdom", "prospecting", "inventory_space"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %#v, want %#v", got, want)
	}
}

func TestNormalizePrioritiesPreservesExplicitInventoryPosition(t *testing.T) {
	priorities := []string{"wisdom", "inventory_space", "prospecting"}
	got, err := NormalizePriorities(priorities)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"wisdom", "inventory_space", "prospecting"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("priorities = %#v, want %#v", got, want)
	}
}

func TestNormalizePrioritiesRejectsInvalidEffect(t *testing.T) {
	_, err := NormalizePriorities([]string{"wisdom", "unknown"})
	if err == nil {
		t.Fatal("expected invalid effect error")
	}
}

func TestNormalizePrioritiesRejectsDuplicateEffect(t *testing.T) {
	_, err := NormalizePriorities([]string{"wisdom", "wisdom"})
	if err == nil {
		t.Fatal("expected duplicate effect error")
	}
}

func TestNormalizePrioritiesRejectsMultipleGatheringSkills(t *testing.T) {
	_, err := NormalizePriorities([]string{"mining", "woodcutting"})
	if err == nil {
		t.Fatal("expected multiple gathering skills error")
	}
}
