package database

import (
	_ "embed"

	"github.com/br-lemes/golem/pkg/schemas"
)

//go:embed npcs_details.json
var npcsDetails []byte

var NpcsDetails = newStore(jsonLoader[schemas.NPCSchema](npcsDetails), func(item *schemas.NPCSchema) string { return item.Code })

var NpcsItems = newStore(func() []schemas.NPCItemSchema {
	npcs := NpcsDetails.All()
	var items []schemas.NPCItemSchema
	for _, npc := range npcs {
		if npc.Items == nil {
			//+gocover:ignore:block all NPCs include item data
			continue
		}
		for _, item := range *npc.Items {
			items = append(items, schemas.NPCItemSchema{
				BuyPrice:  item.BuyPrice,
				Code:      item.Code,
				Currency:  item.Currency,
				Npc:       npc.Code,
				SellPrice: item.SellPrice,
			})
		}
	}
	return items
}, func(item *schemas.NPCItemSchema) string { return item.Code })
