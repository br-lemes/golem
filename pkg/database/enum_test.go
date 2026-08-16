package database

import "testing"

func TestEnumsAreNonEmptyAndUnique(t *testing.T) {
	names := EnumNames()
	if len(names) == 0 {
		t.Fatal("enum names are empty")
	}

	seenNames := make(map[string]struct{}, len(names))
	for _, name := range names {
		if name == "" {
			t.Fatal("enum names contain an empty value")
		}
		_, exists := seenNames[name]
		if exists {
			t.Fatalf("enum names contain duplicate value %q", name)
		}
		seenNames[name] = struct{}{}

		values := Enum(name)
		if len(values) == 0 {
			t.Fatalf("enum %q is empty", name)
		}

		seenValues := make(map[string]struct{}, len(values))
		for _, value := range values {
			if value == "" {
				t.Fatalf("enum %q contains an empty value", name)
			}
			_, exists := seenValues[value]
			if exists {
				t.Fatalf("enum %q contains duplicate value %q", name, value)
			}
			seenValues[value] = struct{}{}
		}
	}
}

func TestEnumMissingName(t *testing.T) {
	values := Enum("missing")
	if values != nil {
		t.Fatalf("Enum(missing) = %v, want nil", values)
	}
}
