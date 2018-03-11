package main

type Ledger struct {
	priceApi PriceAPI
}

func NewLedger(priceApi PriceAPI) (*Ledger) {
	return &Ledger{priceApi: priceApi}
}

func (l *Ledger) HandleMutations(debit []InventoryMutation, credit []InventoryMutation) ([]LedgerMutation, error) {
	var lml []LedgerMutation

	// Find all unique typeIds and their prices
	typeIds := make(map[int]struct{})
	for _, mut := range append(debit, credit...) {
		if _, exists := typeIds[mut.TypeId]; !exists {
			typeIds[mut.TypeId] = struct{}{}
		}
	}

	uniqueTypeIds := []int{}
	for key := range typeIds {
		uniqueTypeIds = append(uniqueTypeIds, key)
	}

	typeIdsWithPrices, err := l.priceApi.FetchPrices(uniqueTypeIds)
	if err != nil {
		return nil, err
	}

	for _, mut := range append(debit, credit...) {
		price := typeIdsWithPrices[mut.TypeId]
		lml = append(lml, LedgerMutation{
			TypePrice:  price,
			TypeId:     mut.TypeId,
			Change:     price * float32(mut.Change),
			PlayerName: mut.PlayerName,
		})
	}

	return lml, nil
}
