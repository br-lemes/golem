package utils

import (
	"testing"

	"github.com/spf13/cobra"
)

type testFlags struct {
	Name  string   `flag:"name" default:"default-name"`
	Count int      `flag:"count"`
	Quiet bool     `flag:"quiet"`
	Items []string `flag:"item" shorthand:"i" default:"one,two"`
}

type allShorthandFlags struct {
	Name  string   `flag:"name" shorthand:"n"`
	Count int      `flag:"count" shorthand:"c"`
	Quiet bool     `flag:"quiet" shorthand:"q"`
	Items []string `flag:"item"`
}

type unsupportedFlags struct {
	Value float64 `flag:"value"`
}

type unsupportedSliceFlags struct {
	Values []int `flag:"value"`
}

type taggedUnexportedFlags struct {
	value string `flag:"value"`
}

type untaggedFlags struct {
	Value string `flag:"value"`
	Other string
}

func TestFlagsRegisterAndRead(t *testing.T) {
	command := &cobra.Command{Use: "test"}
	err := RegisterFlags[testFlags](command)
	if err != nil {
		t.Fatal(err)
	}
	command.SetArgs([]string{"--count", "3", "--quiet", "--item", "three"})
	err = command.Execute()
	if err != nil {
		t.Fatal(err)
	}

	result, err := ReadFlags[testFlags](command)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "default-name" || result.Count != 3 || !result.Quiet {
		t.Fatalf("unexpected flags: %#v", result)
	}
	if len(result.Items) != 1 || result.Items[0] != "three" {
		t.Fatalf("unexpected items: %#v", result.Items)
	}
}

func TestFlagsUseZeroValueWhenDefaultIsOmitted(t *testing.T) {
	command := &cobra.Command{Use: "test"}
	err := RegisterFlags[testFlags](command)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ReadFlags[testFlags](command)
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 0 || result.Quiet || result.Items == nil {
		t.Fatalf("unexpected zero-value flags: %#v", result)
	}
}

func TestFlagsRegisterAllShorthandVariants(t *testing.T) {
	command := &cobra.Command{Use: "test"}
	err := RegisterFlags[allShorthandFlags](command)
	if err != nil {
		t.Fatal(err)
	}
	command.SetArgs([]string{"-n", "name", "-c", "2", "-q", "--item", "one"})
	err = command.Execute()
	if err != nil {
		t.Fatal(err)
	}
	result, err := ReadFlags[allShorthandFlags](command)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "name" || result.Count != 2 || !result.Quiet || len(result.Items) != 1 {
		t.Fatalf("unexpected flags: %#v", result)
	}
}

func TestFlagsRejectUnsupportedTypes(t *testing.T) {
	err := RegisterFlags[unsupportedFlags](&cobra.Command{})
	if err == nil {
		t.Fatal("RegisterFlags accepted unsupported type")
	}
	err = RegisterFlags[unsupportedSliceFlags](&cobra.Command{})
	if err == nil {
		t.Fatal("RegisterFlags accepted unsupported slice type")
	}
	err = RegisterFlags[int](&cobra.Command{})
	if err == nil {
		t.Fatal("RegisterFlags accepted non-struct type")
	}
}

func TestFlagsRejectInvalidDefaults(t *testing.T) {
	err := RegisterFlags[struct {
		Value bool `flag:"value" default:"invalid"`
	}](&cobra.Command{})
	if err == nil {
		t.Fatal("RegisterFlags accepted invalid bool default")
	}
	err = RegisterFlags[struct {
		Value int `flag:"value" default:"invalid"`
	}](&cobra.Command{})
	if err == nil {
		t.Fatal("RegisterFlags accepted invalid int default")
	}
}

func TestReadFlagsReturnsErrorForMissingFlag(t *testing.T) {
	_, err := ReadFlags[testFlags](&cobra.Command{})
	if err == nil {
		t.Fatal("ReadFlags did not report missing flag")
	}
}

func TestFlagsReadErrorsAndValidDefaults(t *testing.T) {
	_, err := ReadFlags[int](&cobra.Command{})
	if err == nil {
		t.Fatal("ReadFlags accepted non-struct type")
	}
	err = RegisterFlags[struct {
		Value bool `flag:"value" default:"true"`
		Count int  `flag:"count" default:"2"`
	}](&cobra.Command{})
	if err != nil {
		t.Fatal(err)
	}

	command := &cobra.Command{Use: "test"}
	err = RegisterFlags[struct {
		Value bool     `flag:"value"`
		Count int      `flag:"count"`
		Items []string `flag:"item"`
	}](command)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ReadFlags[struct {
		Value bool `flag:"value"`
	}](command)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ReadFlags[struct {
		Count int `flag:"missing"`
	}](command)
	if err == nil {
		t.Fatal("ReadFlags did not report missing integer flag")
	}
	_, err = ReadFlags[struct {
		Items []string `flag:"missing"`
	}](command)
	if err == nil {
		t.Fatal("ReadFlags did not report missing slice flag")
	}
	_, err = ReadFlags[struct {
		Value bool `flag:"missing"`
	}](command)
	if err == nil {
		t.Fatal("ReadFlags did not report missing bool flag")
	}
	_, err = ReadFlags[struct {
		Value string `flag:"value"`
	}](command)
	if err == nil {
		t.Fatal("ReadFlags did not report mismatched flag type")
	}
}

func TestFlagsRejectUnexportedField(t *testing.T) {
	_ = taggedUnexportedFlags{}.value

	_, _, err := flagFields(&taggedUnexportedFlags{})
	if err == nil {
		t.Fatal("flagFields accepted unexported field")
	}
}

func TestFlagsIgnoreUntaggedFields(t *testing.T) {
	_, fields, err := flagFields(&untaggedFlags{})
	if err != nil || len(fields) != 1 {
		t.Fatalf("flagFields() = (%v, %d), want one field", err, len(fields))
	}
}
