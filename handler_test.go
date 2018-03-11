package main

import (
	"testing"
	"net/http"
	"reflect"
)

/**
This acts as a sort of integration test since a lot of the handle functionality in the ledger and inventory classes
is tested by this as well. We add 2 stacks to the inventory (total 2000 and then remove 1500 from it, so we have a
removed and a partial stack left).
 */
func TestHandler_Process(t *testing.T) {

	h := Handler{
		ledger: &Ledger{priceApi: &EveMarketerAPI{client: &http.Client{}}},
		inv: NewInventory(),
	}

	typeOxygen := Type{
		TypeID:   3683,
		TypeName: "Oxygen",
	}

	playerGebbetje := Player{
		Name: "Gebbetje",
	}

	playerSwaffeltje := Player{
		Name: "Swaffeltje",
	}

	playerEpicCyno := Player{
		Name: "EpicCyno",
	}

	invMuts, ledgerMuts, err := h.Process([]Transaction{
		{
			Id:         "1",
			PlayerName: playerGebbetje.Name,
			TypeName:   typeOxygen.TypeName,
			Quantity:   1000,
			Action:     ACTION_LOCK,
		},
		{
			Id:         "2",
			PlayerName: playerEpicCyno.Name,
			TypeName:   typeOxygen.TypeName,
			Quantity:   1000,
			Action:     ACTION_LOCK,
		},
		{
			Id:         "3",
			PlayerName: playerSwaffeltje.Name,
			TypeName:   typeOxygen.TypeName,
			Quantity:   1500,
			Action:     ACTION_UNLOCK,
		},
	})

	if err != nil {
		t.Fatalf("Didn't expect error, but got %s", err)
	}

	// Make sure only EpicCyno has 500 left
	expectedContents := map[int][]InventoryStack{
		typeOxygen.TypeID: {
			{PlayerName: playerEpicCyno.Name, Amount: 500},
		},
	}

	if !reflect.DeepEqual(expectedContents, h.inv.contents) {
		t.Fatalf("Expected %v, but got %v", expectedContents, h.inv.contents)
	}

	if len(ledgerMuts) != 3 {
		t.Fatalf("Expected %d ledger mutations, but got %d", 3, len(ledgerMuts))
	}

	if len(invMuts) != 4 {
		t.Fatalf("Expected %d inventory mutations, but got %d", 4, len(invMuts))
	}
}
