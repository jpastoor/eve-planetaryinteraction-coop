package main

type Handler struct {
	inv           *Inventory
	ledger        *Ledger
	playerFetcher PlayerFetcher
	typeFetcher   TypeFetcher
}

func (h *Handler) Process(ts []Transaction) (invCreditMuts []InventoryMutation, invDebitMuts []InventoryMutation,  err error) {

	data := map[int]map[string]int{}
	for _, t := range ts {

		// First we update the playerName of each transaction
		who, err := h.playerFetcher.getOrCreatePlayerByName(t.PlayerName)
		if err != nil {
			return nil, nil, err
		}

		if who.Main != "" {
			t.PlayerName = who.Main
		}

		if t.MarkedForCorp {
			t.PlayerName = "ADHC"
		}

		// Then we make a sum of all the item locks and unlocks per player
		transType, err := h.typeFetcher.getTypeByName(t.TypeName)
		if err != nil {
			return nil, nil,err
		}

		if _, exists := data[transType.TypeID]; !exists {
			data[transType.TypeID] = map[string]int{}
		}

		if _, exists := data[transType.TypeID][t.PlayerName]; !exists {
			data[transType.TypeID][t.PlayerName] = 0
		}

		if t.Action == ACTION_LOCK {
			data[transType.TypeID][t.PlayerName] += t.Quantity
		}

		if t.Action == ACTION_UNLOCK {
			data[transType.TypeID][t.PlayerName] -= t.Quantity
		}
	}

	var creditMutations []InventoryMutation
	var debitMutations []InventoryMutation

	// We then process all the positive amounts
	for typeId, playerMap := range data {
		for playerName, amount := range playerMap {
			if amount > 0 {
				invMut := InventoryMutation{
					PlayerName: playerName,
					TypeId:     typeId,
					Change:     amount,
				}

				h.inv.Add(invMut)
			}
		}
	}

	// Then we process all the negative amounts
	for typeId, playerMap := range data {
		for playerName, amount := range playerMap {
			if amount < 0 {
				invMut := InventoryMutation{
					PlayerName: playerName,
					TypeId:     typeId,
					Change:     amount,
				}

				nCreditMuts, nDebitMuts := h.inv.Sub(invMut)
				creditMutations = append(creditMutations, nCreditMuts...)
				debitMutations = append(debitMutations, nDebitMuts...)
			}
		}
	}

	return creditMutations, debitMutations, nil
}
