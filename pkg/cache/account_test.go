package cache

import "testing"

func TestAccountCache(t *testing.T) {
	initializeTestCache(t)
	SaveAccount("account-name")
	if GetAccount() != "account-name" {
		t.Fatal("account was not saved")
	}
	CleanAccount()
	if GetAccount() != "" {
		t.Fatal("account was not cleaned")
	}
}
