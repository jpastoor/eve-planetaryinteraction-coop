package main

import "fmt"

type Ledger struct {
	priceApi PriceAPI
}

func (l *Ledger) HandleUnlock(t Transaction, invMuts []InventoryMutation) ([]LedgerMutation, error) {
	var lml []LedgerMutation
	// Fetch the price of the type
	price, err := l.priceApi.FetchPrice(*t.Type)
	if err != nil {
		return lml, fmt.Errorf("could not fetch price of %s: %s", t.Type.Name, err)
	}

	// Debit mutation
	lml = append(lml, LedgerMutation{
		Transaction:   &t,
		TransactionId: t.Id,
		TypePrice:     price,
		Change:        price * float32(t.Quantity),
		PlayerName:    t.PlayerName,
	})

	// Credit mutations
	for _, invMut := range invMuts {
		lml = append(lml, LedgerMutation{
			Transaction:   &t,
			TransactionId: t.Id,
			TypePrice:     price,
			Change:        price * float32(invMut.Change),
			PlayerName:    invMut.PlayerName,
		})
	}

	return lml, nil
}
