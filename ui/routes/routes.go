package routes

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/a-h/templ"
)

type Segment struct {
	Icon  templ.Component
	Label string
	Path  string
}

type Definition struct {
	Order   int
	Label   string
	Path    string
	Icon    templ.Component
	Handler func(*http.Request) (templ.Component, int)
}

type Node struct {
	Children []*Node
	Handler  func(*http.Request) (templ.Component, int)
	Parent   *Node
	Path     string
	Segment  Segment
	order    int
}

type Route struct {
	Label string
	Route *Node
}

var root = &Node{}
var notFound templ.Component

func Navigation() []*Node {
	return root.Children
}

func RegisterRoute(definition Definition) *Node {
	if definition.Handler == nil {
		panic("route handler is required")
	}
	if definition.Path == "" || definition.Path[0] != '/' {
		panic(fmt.Sprintf("invalid route path: %q", definition.Path))
	}
	path := strings.TrimPrefix(definition.Path, "/")
	var paths []string
	if path != "" {
		paths = strings.Split(path, "/")
		for _, segment := range paths {
			if segment == "" {
				panic(fmt.Sprintf("invalid route path: %q", definition.Path))
			}
		}
	}

	node := root
	for _, segmentPath := range paths {
		var child *Node
		for _, candidate := range node.Children {
			if candidate.Segment.Path == segmentPath {
				child = candidate
				break
			}
		}
		if child == nil {
			child = &Node{Parent: node, Segment: Segment{Path: segmentPath}}
			node.Children = append(node.Children, child)
		}
		node = child
	}
	if node.Handler != nil {
		panic(fmt.Sprintf("duplicate page route: %s", definition.Path))
	}
	node.Path = definition.Path
	node.Handler = definition.Handler
	node.order = definition.Order
	if len(paths) > 0 {
		node.Segment.Label = definition.Label
		node.Segment.Icon = definition.Icon
	}
	return node
}

func RegisterNotFound(component templ.Component) {
	if notFound != nil {
		panic("duplicate not found page")
	}
	notFound = component
}

func RegisterRoutes(mux *http.ServeMux) {
	sortChildren(root)
	registerRoutes(mux, root)
}

func sortChildren(node *Node) {
	slices.SortFunc(node.Children, func(left, right *Node) int {
		leftOrdered := left.order > 0
		rightOrdered := right.order > 0
		if leftOrdered != rightOrdered {
			if leftOrdered {
				return -1
			}
			return 1
		}
		if leftOrdered && left.order != right.order {
			return cmp.Compare(left.order, right.order)
		}
		return strings.Compare(left.Segment.Label, right.Segment.Label)
	})
	for _, child := range node.Children {
		sortChildren(child)
	}
}

func registerRoutes(mux *http.ServeMux, node *Node) {
	if node.Handler != nil {
		pattern := "GET " + node.Path
		if node == root {
			pattern = "GET /{$}"
		}
		mux.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			component, status := node.Handler(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
			}
			templ.Handler(component).ServeHTTP(w, r)
		}))
	}
	for _, child := range node.Children {
		registerRoutes(mux, child)
	}
	if node == root && notFound != nil {
		mux.Handle("GET /", templ.Handler(notFound, templ.WithStatus(http.StatusNotFound)))
	}
}

func Trail(node *Node) []*Node {
	var trail []*Node
	for node != nil && node != root {
		trail = append(trail, node)
		node = node.Parent
	}
	for left, right := 0, len(trail)-1; left < right; left, right = left+1, right-1 {
		trail[left], trail[right] = trail[right], trail[left]
	}
	return trail
}
