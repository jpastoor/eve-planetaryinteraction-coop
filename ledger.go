package main

import "fmt"

type Ledger struct {
	priceApi PriceAPI
}

func (l *Ledger) HandleMutations(debit []InventoryMutation, credit []InventoryMutation) ([]LedgerMutation, error) {
	var lml []LedgerMutation

	// Find all unique typeIds and their prices
	typeIds := map[int]float32{}
	for _, mut := range append(debit, credit...) {
		if _, exists := typeIds[mut.TypeId]; !exists {
			price, err := l.priceApi.FetchPrice(Type{TypeId: mut.TypeId})
			if err != nil {
				return lml, fmt.Errorf("could not fetch price of %d: %s", mut.TypeId, err)
			}

			typeIds[mut.TypeId] = price
		}
	}

	for _, mut := range append(debit, credit...) {
		price := typeIds[mut.TypeId]
		lml = append(lml, LedgerMutation{
			TypePrice:  price,
			Change:     price * float32(mut.Change),
			PlayerName: mut.PlayerName,
		})
	}

	return lml, nil
}
