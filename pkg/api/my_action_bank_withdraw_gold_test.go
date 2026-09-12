package api

import (
	"net/http"
	"testing"
)

func TestMyActionBankWithdrawGold(t *testing.T) {
	client := bankActionClient(t, "/my/hero/action/bank/withdraw/gold", []byte(`{"quantity":5}`), []byte(`{"data":{"bank":{},"character":{}}}`))
	_, err := client.MyActionBankWithdrawGold("hero", 5)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMyActionBankWithdrawGoldErrors(t *testing.T) {
	for name, client := range map[string]*Client{
		"request error": newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"action failed"}}`))),
		"JSON error":    newTestClient(responseTransport(http.StatusOK, []byte("invalid json"))),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := client.MyActionBankWithdrawGold("hero", 1)
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
