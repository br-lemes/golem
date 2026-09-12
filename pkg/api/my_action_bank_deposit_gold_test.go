package api

import (
	"net/http"
	"testing"
)

func TestMyActionBankDepositGold(t *testing.T) {
	client := bankActionClient(t, "/my/hero/action/bank/deposit/gold", []byte(`{"quantity":5}`), []byte(`{"data":{"bank":{},"character":{}}}`))
	_, err := client.MyActionBankDepositGold("hero", 5)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMyActionBankDepositGoldErrors(t *testing.T) {
	for name, client := range map[string]*Client{
		"request error": newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"action failed"}}`))),
		"JSON error":    newTestClient(responseTransport(http.StatusOK, []byte("invalid json"))),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := client.MyActionBankDepositGold("hero", 1)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
