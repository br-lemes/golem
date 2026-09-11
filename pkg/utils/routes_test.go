package utils

import "testing"

func TestRoutesAndRouteData(t *testing.T) {
	routes, err := GetRoutes()
	if err != nil || len(routes) == 0 {
		t.Fatalf("GetRoutes() = %d routes, %v", len(routes), err)
	}
	for i := 1; i < len(routes); i++ {
		if routePath(routes[i-1]) > routePath(routes[i]) {
			t.Fatal("GetRoutes returned unsorted paths")
		}
	}
	if len(GetRoutesCompletion()) != len(routes) {
		t.Fatal("GetRoutesCompletion returned the wrong number of routes")
	}
	route, err := GetRoute("/items")
	if err != nil || len(route) == 0 {
		t.Fatalf("GetRoute(/items) = %#v, %v", route, err)
	}
	_, err = GetRoute("/does-not-exist")
	if err == nil {
		t.Fatal("GetRoute accepted an unknown path")
	}
}

func TestGetRouteExtractsObjectRequestBody(t *testing.T) {
	route, err := GetRoute("/my/{name}/action/use")
	if err != nil {
		t.Fatal(err)
	}
	for _, details := range route {
		if len(details.RequestBody) == 0 {
			t.Fatalf("GetRoute() returned no request body: %#v", details)
		}
		for _, field := range details.RequestBody {
			if field.Name == "code" && field.Type == "string" {
				return
			}
		}
	}
	t.Fatalf("GetRoute() did not extract expected body field: %#v", route)
}

func routePath(route map[string]string) string {
	for _, path := range route {
		return path
	}
	return ""
}
