package api

import (
	"net/http"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestMyActionBankWithdrawItem(t *testing.T) {
	path := "/my/hero/action/bank/withdraw/item"
	request := []byte(`[{"code":"iron","quantity":2}]`)
	response := []byte(`{"data":{"bank":[],"character":{}}}`)
	client := bankActionClient(t, path, request, response)
	items := []schemas.SimpleItemSchema{{Code: "iron", Quantity: 2}}
	_, err := client.MyActionBankWithdrawItem("hero", items)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMyActionBankWithdrawItemErrors(t *testing.T) {
	for name, client := range map[string]*Client{
		"request error": newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"action failed"}}`))),
		"JSON error":    newTestClient(responseTransport(http.StatusOK, []byte("invalid json"))),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := client.MyActionBankWithdrawItem("hero", nil)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
