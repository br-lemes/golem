package completion

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionBuilderUsesFirstMatchingValidator(t *testing.T) {
	builder := Custom(1, func() []string { return []string{"first"} })
	builder.Custom(1, func() []string { return []string{"second"} })
	complete := builder.Build()
	got, directive := complete(&cobra.Command{}, nil, "")
	if directive != cobra.ShellCompDirectiveNoFileComp || !reflect.DeepEqual(got, []string{"first"}) {
		t.Fatalf("completion = %#v, %v", got, directive)
	}
}

func TestCompletionBuilderSkipsConsumedArguments(t *testing.T) {
	builder := Custom(1, func() []string { return []string{"first"} })
	builder.Custom(1, func() []string { return []string{"second"} })
	complete := builder.Build()
	got, _ := complete(&cobra.Command{}, []string{"used"}, "")
	if !reflect.DeepEqual(got, []string{"second"}) {
		t.Fatalf("completion after consumed argument = %#v", got)
	}
}

func TestCompletionBuilderSkipsValidatorWithNoConsumption(t *testing.T) {
	seenOffset := -1
	first := func(*cobra.Command, []string, string, int) ([]string, int, bool) {
		return nil, 0, false
	}
	second := func(_ *cobra.Command, args []string, _ string, offset int) ([]string, int, bool) {
		seenOffset = offset
		return []string{"match"}, 1, true
	}
	builder := &CompletionBuilder{}
	builder.validators = []ArgValidator{first, second}
	complete := builder.Build()
	got, _ := complete(&cobra.Command{}, []string{"one", "two"}, "")
	if seenOffset != 2 || !reflect.DeepEqual(got, []string{"match"}) {
		t.Fatalf("completion = %#v, offset = %d", got, seenOffset)
	}
}

func TestCompletionBuilderReturnsNoSuggestionsWhenUnmatched(t *testing.T) {
	validator := func(*cobra.Command, []string, string, int) ([]string, int, bool) {
		return nil, 1, false
	}
	builder := &CompletionBuilder{validators: []ArgValidator{validator}}
	complete := builder.Build()
	got, directive := complete(&cobra.Command{}, nil, "")
	if got != nil || directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("unmatched completion = %#v, %v", got, directive)
	}
}

func TestNoArgs(t *testing.T) {
	got, directive := NoArgs(&cobra.Command{}, nil, "")
	if got != nil || directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("NoArgs() = %#v, %v", got, directive)
	}
}
