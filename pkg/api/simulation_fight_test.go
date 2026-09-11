package api

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/br-lemes/golem/pkg/schemas"
)

func TestSimulationFightSendsRequest(t *testing.T) {
	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = testClient(func(r *http.Request) ([]byte, error) {
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
	})

	got, err := SimulationFight(schemas.CombatSimulationRequestSchema{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, schemas.CombatSimulationDataSchema{}) {
		t.Fatalf("got %#v, want empty result", got)
	}
}

func TestSimulationFightReturnsRequestError(t *testing.T) {
	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"simulation unavailable"}}`))

	_, err := SimulationFight(schemas.CombatSimulationRequestSchema{})
	if err == nil {
		t.Fatal("expected request error")
	}
}

func TestSimulationFightReturnsJSONError(t *testing.T) {
	oldClient := defaultClient
	t.Cleanup(func() { defaultClient = oldClient })
	defaultClient = responseClient(http.StatusOK, []byte("invalid json"))

	_, err := SimulationFight(schemas.CombatSimulationRequestSchema{})
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
