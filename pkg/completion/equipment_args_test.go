package completion

import "testing"

func TestEquipmentArgsIncludesMultiSlotSuffixes(t *testing.T) {
	args := equipmentArgs()
	if len(args) == 0 {
		t.Fatal("equipment arguments returned no values")
	}
	foundSuffix := false
	for _, arg := range args {
		if len(arg) > 2 && arg[len(arg)-2] == '@' {
			foundSuffix = true
			break
		}
	}
	if !foundSuffix {
		t.Fatal("equipment arguments contain no multi-slot suffix")
	}
	suggestions, _ := EquipmentArgs(1).Build()(nil, nil, "")
	if len(suggestions) == 0 {
		t.Fatal("equipment argument builder was not created")
	}
}
