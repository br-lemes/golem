package cmd

import (
	"strings"
	"testing"
)

func TestParseCliEquipment(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errContains string
	}{
		{
			name:        "Valid item with single slot no indicator",
			input:       "iron_shield",
			expectError: false,
		},
		{
			name:        "Valid multi-slot item with correct slot",
			input:       "copper_ring@1",
			expectError: false,
		},
		{
			name:        "Valid utility slot with quantity",
			input:       "health_potion@1x5",
			expectError: false,
		},
		{
			name:        "Invalid slot zero provided for multi slot",
			input:       "copper_ring@0",
			expectError: true,
			errContains: "invalid slot number",
		},
		{
			name:        "Invalid slot zero provided for single slot",
			input:       "iron_shield@0",
			expectError: true,
			errContains: "invalid slot number",
		},
		{
			name:        "Invalid negative slot provided",
			input:       "copper_ring@-2",
			expectError: true,
			errContains: "invalid slot number",
		},
		{
			name:        "Slot number higher than available slots",
			input:       "copper_ring@3",
			expectError: true,
			errContains: "invalid slot number",
		},
		{
			name:        "Multi-slot item missing slot indicator",
			input:       "copper_ring",
			expectError: true,
			errContains: "must specify slot",
		},
		{
			name:        "Single slot item trying to specify a slot number",
			input:       "iron_shield@1",
			expectError: true,
			errContains: "cannot specify slot number for item with single slot",
		},
		{
			name:        "Invalid quantity for non-utility single slot",
			input:       "iron_shield@1x5",
			expectError: true,
			errContains: "cannot specify slot number for item with single slot",
		},
		{
			name:        "Invalid quantity for non-utility multi slot",
			input:       "copper_ring@1x5",
			expectError: true,
			errContains: "cannot specify quantity for non-utility slot",
		},
		{
			name:        "Non-equipment item",
			input:       "apple",
			expectError: true,
			errContains: "item type cannot be equipped",
		},
		{
			name:        "Non-existent item",
			input:       "invalid_item",
			expectError: true,
			errContains: "item not found in database",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseCliEquipment(tc.input)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil",
						tc.errContains)
				}
				if !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("expected error %q to contain %q", err.Error(),
						tc.errContains)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			}
		})
	}
}
