package api

import (
	"net/http"
	"testing"
)

func TestGetBuildsQuery(t *testing.T) {
	tests := []struct {
		name string
		data map[string]string
		want string
	}{
		{
			name: "nil data",
			want: "/items",
		},
		{
			name: "empty data",
			data: map[string]string{},
			want: "/items",
		},
		{
			name: "empty values",
			data: map[string]string{"name": "", "type": ""},
			want: "/items",
		},
		{
			name: "encoded values",
			data: map[string]string{"name": "iron ore", "type": "resource"},
			want: "/items?name=iron+ore&type=resource",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			oldClient := defaultClient
			t.Cleanup(func() { defaultClient = oldClient })
			defaultClient = testClient(func(r *http.Request) ([]byte, error) {
				if r.Method != http.MethodGet {
					t.Errorf("method = %s, want GET", r.Method)
				}
				if r.URL.RequestURI() != test.want {
					t.Errorf("request URI = %s, want %s", r.URL.RequestURI(), test.want)
				}
				return []byte(`result`), nil
			})

			got, err := Get("/items", test.data)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != "result" {
				t.Fatalf("response = %q, want result", got)
			}
		})
	}
}
