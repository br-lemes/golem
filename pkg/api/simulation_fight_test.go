package api

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSimulationFightSendsRequest(t *testing.T) {
	client := newTestClient(testTransport(func(r *http.Request) ([]byte, error) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/simulation/fight" {
			t.Errorf("path = %s, want /simulation/fight", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		wantBody := `{"characters":null,"iterations":0,"monster":""}`
		if string(body) != wantBody {
			t.Errorf("body = %s, want %s", body, wantBody)
		}
		return []byte(`{"data":{}}`), nil
	}))

	got, err := client.SimulationFight(schemas.CombatSimulationRequestSchema{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, schemas.CombatSimulationDataSchema{}) {
		t.Fatalf("got %#v, want empty result", got)
	}
}

func TestSimulationFightReturnsRequestError(t *testing.T) {
	client := newTestClient(responseTransport(http.StatusBadRequest, []byte(`{"error":{"message":"simulation unavailable"}}`)))

	_, err := client.SimulationFight(schemas.CombatSimulationRequestSchema{})
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestSimulationFightReturnsJSONError(t *testing.T) {
	client := newTestClient(responseTransport(http.StatusOK, []byte("invalid json")))

	_, err := client.SimulationFight(schemas.CombatSimulationRequestSchema{})
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
