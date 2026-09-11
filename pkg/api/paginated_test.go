package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/br-lemes/golem/pkg/cache"
)

func TestPaginatedEndpoints(t *testing.T) {
	tests := []struct {
		name string
		path string
		size int
		call func() (int, error)
	}{
		{
			name: "items",
			path: "/items",
			size: ItemsSize,
			call: func() (int, error) {
				items, err := Items(ItemsOptions{})
				return len(items), err
			},
		},
		{
			name: "monsters",
			path: "/monsters",
			size: MonstersSize,
			call: func() (int, error) {
				items, err := Monsters(MonstersOptions{})
				return len(items), err
			},
		},
		{
			name: "maps",
			path: "/maps",
			size: MapsSize,
			call: func() (int, error) {
				items, err := Maps(MapsOptions{})
				return len(items), err
			},
		},
		{
			name: "effects",
			path: "/effects",
			size: EffectsSize,
			call: func() (int, error) {
				items, err := Effects(EffectsOptions{})
				return len(items), err
			},
		},
		{
			name: "resources",
			path: "/resources",
			size: ResourcesSize,
			call: func() (int, error) {
				items, err := Resources(ResourcesOptions{})
				return len(items), err
			},
		},
		{
			name: "events",
			path: "/events",
			size: EventsSize,
			call: func() (int, error) {
				items, err := Events(EventsOptions{})
				return len(items), err
			},
		},
		{
			name: "achievements",
			path: "/achievements",
			size: AchievementsSize,
			call: func() (int, error) {
				items, err := Achievements(AchievementsOptions{})
				return len(items), err
			},
		},
		{
			name: "tasks",
			path: "/tasks/list",
			size: TasksListSize,
			call: func() (int, error) {
				items, err := TasksList(TasksListOptions{})
				return len(items), err
			},
		},
		{
			name: "npcs",
			path: "/npcs/details",
			size: NpcsDetailsSize,
			call: func() (int, error) {
				items, err := NpcsDetails(NpcsDetailsOptions{})
				return len(items), err
			},
		},
		{
			name: "grandexchange orders",
			path: "/grandexchange/orders",
			size: GrandexchangeOrdersSize,
			call: func() (int, error) {
				items, err := GrandexchangeOrders(GrandexchangeOrdersOptions{})
				return len(items), err
			},
		},
		{
			name: "my grandexchange orders",
			path: "/my/grandexchange/orders",
			size: MyGrandexchangeOrdersSize,
			call: func() (int, error) {
				items, err := MyGrandexchangeOrders(MyGrandexchangeOrdersOptions{})
				return len(items), err
			},
		},
		{
			name: "my bank items",
			path: "/my/bank/items",
			size: MyBankItemsSize,
			call: func() (int, error) {
				items, err := MyBankItems()
				return len(items), err
			},
		},
		{
			name: "active events",
			path: "/events/active",
			size: EventsActiveSize,
			call: func() (int, error) {
				items, err := EventsActive()
				return len(items), err
			},
		},
		{
			name: "account achievements",
			path: "/accounts/test/achievements",
			size: AccountsAchievementsSize,
			call: func() (int, error) {
				items, err := AccountsAchievements("test")
				return len(items), err
			},
		},
		{
			name: "grandexchange history",
			path: "/grandexchange/history/iron",
			size: GrandexchangeHistorySize,
			call: func() (int, error) {
				items, err := GrandexchangeHistory("iron", GrandexchangeHistoryOptions{})
				return len(items), err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache.CleanBankItems()
			oldBaseURL := baseURL
			oldClient := defaultClient
			t.Cleanup(func() {
				baseURL = oldBaseURL
				defaultClient = oldClient
			})

			requests := 0
			defaultClient = testClient(func(r *http.Request) ([]byte, error) {
				requests++
				if r.URL.Path != test.path {
					return nil, fmt.Errorf("path = %s, want %s", r.URL.Path, test.path)
				}
				queryPage(t, r.URL.Query(), requests, test.size)
				return []byte(`{"data":[{},{}],"pages":2}`), nil
			})

			got, err := test.call()
			if err != nil {
				t.Fatal(err)
			}
			if got != 4 {
				t.Fatalf("got %d items, want 4", got)
			}
			if requests != 2 {
				t.Fatalf("got %d requests, want 2", requests)
			}
		})
	}

	for _, test := range tests {
		t.Run(test.name+" returns request error", func(t *testing.T) {
			cache.CleanBankItems()
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = responseClient(http.StatusBadRequest, []byte(`{"error":{"message":"bad request"}}`))
			_, err := test.call()
			if err == nil {
				t.Fatal("expected request error")
			}
		})
	}

	for _, test := range tests {
		t.Run(test.name+" returns JSON error", func(t *testing.T) {
			cache.CleanBankItems()
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = responseClient(http.StatusOK, []byte("invalid json"))
			_, err := test.call()
			if err == nil {
				t.Fatal("expected JSON error")
			}
		})
	}
}

func queryPage(t *testing.T, values url.Values, page, size int) {
	t.Helper()
	if values.Get("page") != fmt.Sprint(page) {
		t.Errorf("page = %q, want %d", values.Get("page"), page)
	}
	if values.Get("size") != fmt.Sprint(size) {
		t.Errorf("size = %q, want %d", values.Get("size"), size)
	}
}

func testClient(handler func(*http.Request) ([]byte, error)) *http.Client {
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body, err := handler(r)
		if err != nil {
			return nil, err
		}
		return jsonResponse(body), nil
	})
	return &http.Client{Transport: transport}
}

func jsonResponse(body []byte) *http.Response {
	return responseClientResponse(http.StatusOK, body)
}

func responseClient(status int, body []byte) *http.Client {
	transport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		return responseClientResponse(status, body), nil
	})
	return &http.Client{Transport: transport}
}

func responseClientResponse(status int, body []byte) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		StatusCode: status,
	}
}
