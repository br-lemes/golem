package cmd

import "testing"

func TestDeleteValidate(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		flags           deleteFlags
		quantityChanged bool
		valid           bool
	}{
		{
			name:  "discard",
			args:  []string{"hero"},
			flags: deleteFlags{Discard: true},
			valid: true,
		},
		{
			name:  "discard with code",
			args:  []string{"hero", "raw_chicken"},
			flags: deleteFlags{Discard: true},
		},
		{
			name:  "discard with all",
			args:  []string{"hero"},
			flags: deleteFlags{Discard: true, All: true},
		},
		{
			name:            "discard with quantity",
			args:            []string{"hero"},
			flags:           deleteFlags{Discard: true},
			quantityChanged: true,
		},
		{
			name: "missing code",
			args: []string{"hero"},
		},
		{
			name:  "all",
			args:  []string{"hero", "raw_chicken"},
			flags: deleteFlags{All: true},
			valid: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := deleteValidate(test.args, test.flags, test.quantityChanged)
			if (err == nil) != test.valid {
				t.Fatalf("deleteValidate() error = %v", err)
			}
		})
	}
}
