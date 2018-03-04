package main

type Handler struct {
	inv    *Inventory
	ledger *Ledger
}

/**
TODO How to make sure transactions don't get double processed
 */

func (h *Handler) Process(ts []Transaction) (invMuts []InventoryMutation, ledgerMuts []LedgerMutation, err error) {
	for _, t := range ts {
		if t.Action == ACTION_UNLOCK {
			// First we update the inventory
			ninvMuts, err := h.inv.handleActionUnlock(t)
			if err != nil {
				return nil, nil, err
			}

			// Based on inventory changes, we know how to update the ledger
			nledgerMuts, err := h.ledger.HandleUnlock(t, invMuts)
			if err != nil {
				return nil, nil, err
			}

			invMuts = append(invMuts, ninvMuts...)
			ledgerMuts = append(ledgerMuts, nledgerMuts...)
		}

		if t.Action == ACTION_LOCK {
			// Update inventory, no change to ledger
			ninvMuts, err := h.inv.handleActionLock(t)
			if err != nil {
				return nil, nil, err
			}

			invMuts = append(invMuts, ninvMuts...)
		}
	}

	return invMuts, ledgerMuts, nil
}
