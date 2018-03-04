package main

import (
	"fmt"
)

/**
Idea here is that you can always rebuild the inventory contents by loading in all the inventory mutations.
So the inventory contents is just an in-memory helper but the real deal is in the mutations that are persisted.

TODO Write tests that ensure contents is in sync with the mutations at all times

 */
type Inventory struct {
	contents map[int][]InventoryStack
}

type InventoryStack struct {
	PlayerName string
	Amount     int
}

type InventoryMutation struct {
	TypeId        int
	TransactionId string
	PlayerName    string
	Change        int
}

func (i *Inventory) Load(muts []InventoryMutation) {
	for _, im := range muts {
		if _, exists := i.contents[im.TypeId]; !exists {
			i.contents[im.TypeId] = []InventoryStack{}
		}

		if im.Change > 0 {
			i.contents[im.TypeId] = append(i.contents[im.TypeId], InventoryStack{
				PlayerName: im.PlayerName,
				Amount:     im.Change,
			})
		}

		// Deduct from existing stacks. and remove empty stacks
		if im.Change < 0 {
			amountLeft := im.Change * -1

			stacks := i.contents[im.TypeId]
			newStacks := []InventoryStack{}

			for _, stack := range stacks {
				if amountLeft > 0 && stack.PlayerName == im.PlayerName {
					var fetchAmount int
					// Stack is bigger than what we need, we do a partial take
					if stack.Amount > amountLeft {
						fetchAmount = amountLeft
					} else {
						// Stack is equal or smaller than that we need, we take all
						fetchAmount = stack.Amount
					}

					stack.Amount -= fetchAmount
					amountLeft -= fetchAmount
				}

				// Only keep the stack if there is anything in it
				if stack.Amount > 0 {
					newStacks = append(newStacks, stack)
				}
			}

			i.contents[im.TypeId] = newStacks
		}
	}
}

func (i *Inventory) handleActionLock(t Transaction) ([]InventoryMutation, error) {
	im := InventoryMutation{
		TypeId:        t.Type.TypeId,
		PlayerName:    t.PlayerName,
		TransactionId: t.Id,
		Change:        t.Quantity,
	}

	invMuts := []InventoryMutation{im}
	i.Load(invMuts)

	return invMuts, nil
}

func (i *Inventory) handleActionUnlock(t Transaction) ([]InventoryMutation, error) {

	amountLeft := t.Quantity
	var muts []InventoryMutation

	for amountLeft > 0 {

		// Find someone with this type in the inventory
		stacks := i.contents[t.Type.TypeId]

		if len(stacks) == 0 {
			return muts, fmt.Errorf("no stacks left, but still need %d of %s for %s", amountLeft, t.TypeName, t.Id)
		}

		stack := stacks[0]
		var fetchAmount int

		// Stack is bigger than what we need, we do a partial take
		if stack.Amount > amountLeft {
			fetchAmount = amountLeft
		} else {
			// Stack is equal or smaller than that we need, we take all
			fetchAmount = stack.Amount
		}

		im := InventoryMutation{
			Change:        fetchAmount * -1,
			TypeId:        t.Type.TypeId,
			PlayerName:    stack.PlayerName,
			TransactionId: t.Id,
		}

		amountLeft -= fetchAmount

		i.Load([]InventoryMutation{im})
		muts = append(muts, im)
	}

	return muts, nil
}
