package api

import (
	"bytes"
	"testing"

	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestPrintDropSchema(t *testing.T) {
	tests := []struct {
		name  string
		drops []schemas.DropSchema
		want  string
	}{
		{
			name: "empty",
			want: "",
		},
		{
			name: "multiple drops",
			drops: []schemas.DropSchema{
				{Code: "iron", Quantity: 2},
				{Code: "wood", Quantity: 1},
			},
			want: ", Drops: 2x iron, 1x wood",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldStdout := console.Stdout
			t.Cleanup(func() { console.Stdout = oldStdout })
			output := &bytes.Buffer{}
			console.Stdout = output

			printDropSchema(test.drops)

			if output.String() != test.want {
				t.Errorf("output = %q, want %q", output.String(), test.want)
			}
		})
	}
}
