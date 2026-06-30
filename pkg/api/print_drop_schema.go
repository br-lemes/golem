package api

import (
	"github.com/br-lemes/golem/pkg/console"
	"github.com/br-lemes/golem/pkg/schemas"
)

func printDropSchema(drops []schemas.DropSchema) {
	length := len(drops)
	if length == 0 {
		return
	}
	console.Printf(", Drops: ")
	for i, drop := range drops {
		console.Printf("%dx %s", drop.Quantity, drop.Code)

		if i < length-1 {
			console.Printf(", ")
		}
	}
}
