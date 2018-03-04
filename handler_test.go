package main

import (
	"testing"
	"net/http"
	"reflect"
)

func TestHandler_Process(t *testing.T) {

	h := Handler{
		ledger: &Ledger{priceApi: &EveMarketerAPI{client: &http.Client{}}},
		inv: &Inventory{
			contents: make(map[int][]InventoryStack),
		},
	}

	typeOxygen := Type{
		TypeId: 3683,
		Name:   "Oxygen",
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
			Who:        &playerGebbetje,
			TypeName:   typeOxygen.Name,
			Type:       &typeOxygen,
			Quantity:   1000,
			Action:     ACTION_LOCK,
		},
		{
			Id:         "2",
			PlayerName: playerEpicCyno.Name,
			Who:        &playerEpicCyno,
			TypeName:   typeOxygen.Name,
			Type:       &typeOxygen,
			Quantity:   1000,
			Action:     ACTION_LOCK,
		},
		{
			Id:         "3",
			PlayerName: playerSwaffeltje.Name,
			Who:        &playerSwaffeltje,
			TypeName:   typeOxygen.Name,
			Type:       &typeOxygen,
			Quantity:   1500,
			Action:     ACTION_UNLOCK,
		},
	})

	if err != nil {
		t.Fatalf("Didn't expect error, but got %s", err)
	}

	// Make sure only EpicCyno has 500 left
	expectedContents := map[int][]InventoryStack{
		typeOxygen.TypeId: {
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