package cmd

import "testing"

func TestStockShowCommand(t *testing.T) {
	command, _, err := stockCmd.Find([]string{"show"})
	if err != nil || command != stockShowCmd || stockCmd.Flags().Lookup("target") == nil || stockShowCmd.Flags().Lookup("target") == nil {
		t.Fatalf("stock show command = %#v, %v", command, err)
	}
}
