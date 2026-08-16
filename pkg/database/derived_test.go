package database

import "testing"

func TestDerivedStringListsAreNonEmptyAndUnique(t *testing.T) {
	tests := []struct {
		name string
		list func() []string
	}{
		{name: "item types", list: ItemTypes},
		{name: "event content codes", list: EventContentCodes},
		{name: "task skills", list: TaskSkills},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values := test.list()
			if len(values) == 0 {
				t.Fatal("derived list is empty")
			}

			seen := make(map[string]struct{}, len(values))
			for _, value := range values {
				if value == "" {
					t.Fatal("derived list contains an empty value")
				}
				_, exists := seen[value]
				if exists {
					t.Fatalf("derived list contains duplicate value %q", value)
				}
				seen[value] = struct{}{}
			}
		})
	}
}
