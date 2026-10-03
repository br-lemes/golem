package routes

import "testing"

func TestSortChildren(t *testing.T) {
	root := &Node{
		Children: []*Node{
			{Segment: Segment{Label: "Zebra"}},
			{Segment: Segment{Label: "Alpha"}},
			{Segment: Segment{Label: "Beta"}, order: -1},
			{Segment: Segment{Label: "Gamma"}, order: 1},
			{Segment: Segment{Label: "Delta"}, order: 2},
		},
	}
	sortChildren(root)

	want := []string{"Gamma", "Delta", "Alpha", "Beta", "Zebra"}
	for index, label := range want {
		if root.Children[index].Segment.Label != label {
			t.Fatalf("children[%d] = %q, want %q", index, root.Children[index].Segment.Label, label)
		}
	}
}
