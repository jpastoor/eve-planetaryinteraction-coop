package main

/**
Idea here is that you can always rebuild the inventory contents by loading in all the inventory mutations.
So the inventory contents is just an in-memory helper but the real deal is in the mutations that are persisted.

TODO Write tests that ensure contents is in sync with the mutations at all times

 */
type Inventory struct {
	contents map[int][]InventoryStack
}

func NewInventory() (*Inventory) {
	return &Inventory{
		contents: make(map[int][]InventoryStack),
	}
}

type InventoryStack struct {
	PlayerName string
	Amount     int
}

type InventoryMutation struct {
	TypeId     int
	PlayerName string
	TypeName   string
	Change     int
}

func (inv *Inventory) Add(im InventoryMutation) {

	if _, exists := inv.contents[im.TypeId]; !exists {
		inv.contents[im.TypeId] = []InventoryStack{}
	}

	inv.contents[im.TypeId] = append(inv.contents[im.TypeId], InventoryStack{
		PlayerName: im.PlayerName,
		Amount:     im.Change,
	})
}

func (inv *Inventory) Sub(im InventoryMutation) ([]InventoryMutation) {
	amountLeft := im.Change * -1

	stacks := inv.contents[im.TypeId]
	var newStacks []InventoryStack

	var mutationsForLedger []InventoryMutation

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

			mutationsForLedger = append(mutationsForLedger, InventoryMutation{
				Change:     fetchAmount,
				TypeId:     im.TypeId,
				PlayerName: stack.PlayerName,
			})
		}

		// Only keep the stack if there is anything in it
		if stack.Amount > 0 {
			newStacks = append(newStacks, stack)
		}
	}

	inv.contents[im.TypeId] = newStacks

	return mutationsForLedger
}
