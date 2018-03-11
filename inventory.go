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

	// TODO Read from the InventoryState table and prefill

	return &Inventory{
		contents: make(map[int][]InventoryStack),
	}
}

type InventoryStack struct {
	PlayerName string
	Amount     int
}

/**
This type will be persisted in the database at every commit
 */
type InventoryState struct {
	CommitId   int    `gorm:"primary_key auto_increment=false"`
	TypeId     int    `gorm:"primary_key auto_increment=false"`
	PlayerName string `gorm:"primary_key auto_increment=false"`
	Amount     int
}

type InventoryMutation struct {
	TypeId     int
	PlayerName string
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

func (inv *Inventory) Sub(im InventoryMutation) (creditMuts []InventoryMutation, debitMuts []InventoryMutation) {
	amountLeft := im.Change * -1

	stacks := inv.contents[im.TypeId]
	var newStacks []InventoryStack

	for _, stack := range stacks {
		if amountLeft > 0 {
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

			creditMuts = append(creditMuts, InventoryMutation{
				Change:     fetchAmount,
				TypeId:     im.TypeId,
				PlayerName: stack.PlayerName,
			})

			debitMuts = append(debitMuts, InventoryMutation{
				Change:     fetchAmount * -1,
				TypeId:     im.TypeId,
				PlayerName: im.PlayerName,
			})
		}

		// Only keep the stack if there is anything in it
		if stack.Amount > 0 {
			newStacks = append(newStacks, stack)
		}
	}

	inv.contents[im.TypeId] = newStacks

	return creditMuts, debitMuts
}
