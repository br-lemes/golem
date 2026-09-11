package api

import (
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyBankReturnsCachedBank(t *testing.T) {
	cache.CleanBank()
	t.Cleanup(cache.CleanBank)
	bank := schemas.BankSchema{}
	cache.SaveBank(bank)

	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"request not expected"}}`))

	got, err := MyBank()
	if err != nil {
		t.Fatal(err)
	}
	if got != bank {
		t.Fatalf("got cached bank %#v, want %#v", got, bank)
	}
}

func TestMyBankFetchesAndCachesBank(t *testing.T) {
	cache.CleanBank()
	t.Cleanup(cache.CleanBank)

	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = testClient(func(r *http.Request) ([]byte, error) {
		if r.URL.Path != "/my/bank" {
			t.Fatalf("path = %s, want /my/bank", r.URL.Path)
		}
		return []byte(`{"data":{}}`), nil
	})

	got, err := MyBank()
	if err != nil {
		t.Fatal(err)
	}
	if got != (schemas.BankSchema{}) {
		t.Fatalf("got bank %#v, want empty bank", got)
	}
	if cache.GetBank() == nil {
		t.Fatal("MyBank() did not cache the bank")
	}
}

func TestMyBankReturnsRequestError(t *testing.T) {
	cache.CleanBank()
	t.Cleanup(cache.CleanBank)
	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"bank unavailable"}}`))

	_, err := MyBank()
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestMyBankReturnsJSONError(t *testing.T) {
	cache.CleanBank()
	t.Cleanup(cache.CleanBank)
	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusOK, []byte("invalid json"))

	_, err := MyBank()
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
