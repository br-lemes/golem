package cmd

import (
	"strconv"
	"testing"
)

func TestMovementOptionsReadPersistentFlags(t *testing.T) {
	oldAllowGold, err := rootCmd.PersistentFlags().GetBool("allow-gold")
	if err != nil {
		t.Fatal(err)
	}
	oldNoTeleport, err := rootCmd.PersistentFlags().GetBool("no-teleport")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := rootCmd.PersistentFlags().Set("allow-gold", strconv.FormatBool(oldAllowGold))
		if err != nil {
			t.Errorf("restore allow-gold flag: %v", err)
		}
		err = rootCmd.PersistentFlags().Set("no-teleport", strconv.FormatBool(oldNoTeleport))
		if err != nil {
			t.Errorf("restore no-teleport flag: %v", err)
		}
	}()
	err = rootCmd.PersistentFlags().Set("allow-gold", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = rootCmd.PersistentFlags().Set("no-teleport", "true")
	if err != nil {
		t.Fatal(err)
	}

	got, err := movementOptions(fightCmd)
	if err != nil {
		t.Fatal(err)
	}
	if !got.AllowGold || !got.NoTeleport {
		t.Fatalf("movement options = %#v, want both persistent options enabled", got)
	}
}
