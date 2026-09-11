package api

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
	"github.com/br-lemes/golem/pkg/schemas"
)

func TestBankActions(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		body     []byte
		response []byte
		call     func() error
	}{
		{
			name:     "deposit gold",
			path:     "/my/hero/action/bank/deposit/gold",
			body:     []byte(`{"quantity":5}`),
			response: []byte(`{"data":{"bank":{},"character":{}}}`),
			call:     func() error { _, err := MyActionBankDepositGold("hero", 5); return err },
		},
		{
			name:     "withdraw gold",
			path:     "/my/hero/action/bank/withdraw/gold",
			body:     []byte(`{"quantity":5}`),
			response: []byte(`{"data":{"bank":{},"character":{}}}`),
			call:     func() error { _, err := MyActionBankWithdrawGold("hero", 5); return err },
		},
		{
			name:     "deposit item",
			path:     "/my/hero/action/bank/deposit/item",
			body:     []byte(`[{"code":"iron","quantity":2}]`),
			response: []byte(`{"data":{"bank":[],"character":{}}}`),
			call: func() error {
				items := []schemas.SimpleItemSchema{{Code: "iron", Quantity: 2}}
				_, err := MyActionBankDepositItem("hero", items)
				return err
			},
		},
		{
			name:     "withdraw item",
			path:     "/my/hero/action/bank/withdraw/item",
			body:     []byte(`[{"code":"iron","quantity":2}]`),
			response: []byte(`{"data":{"bank":[],"character":{}}}`),
			call: func() error {
				items := []schemas.SimpleItemSchema{{Code: "iron", Quantity: 2}}
				_, err := MyActionBankWithdrawItem("hero", items)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache.CleanBank()
			cache.CleanBankItems()
			cache.CleanCharacters()
			t.Cleanup(cache.CleanBank)
			t.Cleanup(cache.CleanBankItems)
			t.Cleanup(cache.CleanCharacters)
			cache.SaveBank(schemas.BankSchema{})
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = testClient(func(r *http.Request) ([]byte, error) {
				if r.Method != http.MethodPost || r.URL.Path != test.path {
					t.Errorf("request = %s %s, want POST %s", r.Method, r.URL.Path, test.path)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(body, test.body) {
					t.Errorf("body = %s, want %s", body, test.body)
				}
				return test.response, nil
			})
			err := test.call()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBankActionsReturnErrors(t *testing.T) {
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "deposit gold",
			call: func() error { _, err := MyActionBankDepositGold("hero", 1); return err },
		},
		{
			name: "withdraw gold",
			call: func() error { _, err := MyActionBankWithdrawGold("hero", 1); return err },
		},
		{
			name: "deposit item",
			call: func() error { _, err := MyActionBankDepositItem("hero", nil); return err },
		},
		{
			name: "withdraw item",
			call: func() error { _, err := MyActionBankWithdrawItem("hero", nil); return err },
		},
	}
	for _, test := range tests {
		t.Run(test.name+" request error", func(t *testing.T) {
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"action failed"}}`))
			err := test.call()
			if err == nil {
				t.Fatal("expected request error")
			}
		})
		t.Run(test.name+" JSON error", func(t *testing.T) {
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = responseClient(http.StatusOK, []byte("invalid json"))
			err := test.call()
			if err == nil {
				t.Fatal("expected JSON error")
			}
		})
	}
}
