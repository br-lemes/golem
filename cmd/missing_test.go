package cmd

import "testing"

func TestMissingValidate(t *testing.T) {
	tests := []struct {
		name  string
		flags missingFlags
		valid bool
	}{
		{
			name:  "default",
			flags: missingFlags{MinLevel: minGameLevel, MaxLevel: maxGameLevel},
			valid: true,
		},
		{
			name:  "range",
			flags: missingFlags{MinLevel: 20, MaxLevel: 30},
			valid: true,
		},
		{
			name:  "minimum below game range",
			flags: missingFlags{MinLevel: 0, MaxLevel: maxGameLevel},
		},
		{
			name:  "invalid range",
			flags: missingFlags{MinLevel: 30, MaxLevel: 20},
		},
		{
			name: "invalid crafting skill",
			flags: missingFlags{
				MinLevel: minGameLevel,
				MaxLevel: maxGameLevel,
				Skill:    []string{"invalid"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := missingValidate(test.flags)
			if (err == nil) != test.valid {
				t.Fatalf("missingValidate() error = %v", err)
			}
		})
	}
}
