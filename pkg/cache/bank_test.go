package cache

import (
	"testing"

	"github.com/br-lemes/golem/pkg/models"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestBankCache(t *testing.T) {
	initializeTestCache(t)
	SaveBank(schemas.BankSchema{Gold: 100, Slots: 20})
	got := GetBank()
	if got == nil || got.Gold != 100 || got.Slots != 20 {
		t.Fatalf("GetBank() = %#v", got)
	}
	UpdateBankGold(250)
	got = GetBank()
	if got == nil || got.Gold != 250 {
		t.Fatalf("GetBank() after update = %#v", got)
	}
	CleanBank()
	if GetBank() != nil {
		t.Fatal("bank was not cleaned")
	}
}

func TestGetBankRejectsCorruptData(t *testing.T) {
	initializeTestCache(t)
	cache.Create(&models.Cache{Name: "bank", Data: "not-json"})
	if GetBank() != nil {
		t.Fatal("corrupt bank data was accepted")
	}
}

func TestBankItemsCache(t *testing.T) {
	initializeTestCache(t)
	items := []schemas.SimpleItemSchema{{Code: "iron_sword", Quantity: 2}}
	SaveBankItems(items)
	got := GetBankItems()
	if len(got) != 1 || got[0].Code != "iron_sword" || got[0].Quantity != 2 {
		t.Fatalf("GetBankItems() = %#v", got)
	}
	CleanBankItems()
	if GetBankItems() != nil {
		t.Fatal("bank items were not cleaned")
	}
}

func TestGetBankItemsRejectsCorruptData(t *testing.T) {
	initializeTestCache(t)
	cache.Create(&models.Cache{Name: "bankItems", Data: "not-json"})
	if GetBankItems() != nil {
		t.Fatal("corrupt bank items were accepted")
	}
}
