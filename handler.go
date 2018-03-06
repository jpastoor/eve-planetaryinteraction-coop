package main

type Handler struct {
	inv    *Inventory
	ledger *Ledger
}

/**
TODO How to make sure transactions don't get double processed
 */
func (h *Handler) Process(ts []Transaction) (invMuts []InventoryMutation, ledgerMuts []LedgerMutation, err error) {

	// First we update the playerName of each transaction
	for _, t := range ts {
		if t.Who.Main != "" {
			t.PlayerName = t.Who.Main
		}

		if t.MarkedForCorp {
			t.PlayerName = "ADHC"
		}
	}

	// Then we make a sum of all the item locks and unlocks per player
	data := map[int]map[string]int{}
	for _, t := range ts {
		if _, exists := data[t.Type.TypeID]; !exists {
			data[t.Type.TypeID] = map[string]int{}
		}

		if _, exists := data[t.Type.TypeID][t.PlayerName]; !exists {
			data[t.Type.TypeID][t.PlayerName] = 0
		}

		if t.Action == ACTION_LOCK {
			data[t.Type.TypeID][t.PlayerName] += t.Quantity
		}

		if t.Action == ACTION_UNLOCK {
			data[t.Type.TypeID][t.PlayerName] -= t.Quantity
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
				debitMutations = append(debitMutations, invMut)
			}
		}
	}

	// Then we process all the negative amounts
	for typeId, playerMap := range data {
		for playerName, amount := range playerMap {
			if amount < 0 {
				creditMutations = append(creditMutations, h.inv.Sub(InventoryMutation{
					PlayerName: playerName,
					TypeId:     typeId,
					Change:     amount,
				})...)
			}
		}
	}

	return invMuts, ledgerMuts, nil
}
