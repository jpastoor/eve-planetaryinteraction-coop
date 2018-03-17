package main

import (
	"testing"
	"net/http"
	"reflect"
)

func TestHandler_Process(t *testing.T) {

	h := Handler{
		ledger: &Ledger{priceApi: &EveMarketerAPI{client: &http.Client{}}},
		inv:    NewInventory(),
		playerFetcher: &DbPlayerFetcherMock{
			cacheByName: map[string]*Player{
				"TestMain":  {Name: "TestMain"},
				"TestAlt":   {Name: "TestAlt", Main: "TestMain"},
				"TestOther": {Name: "TestOther"},
			},
		},
		typeFetcher: &TypeFetcherMock{
			cacheByName: map[string]*Type{
				"Oxygen": {
					TypeID:   3683,
					TypeName: "Oxygen",
				},
			},
		},
	}

	creditMuts, debitMuts, err := h.Process([]Transaction{
		{
			Id:         "1",
			PlayerName: "TestMain",
			TypeName:   "Oxygen",
			Quantity:   1000,
			Action:     ACTION_LOCK,
		},
		{
			Id:         "2",
			PlayerName: "TestOther",
			TypeName:   "Oxygen",
			Quantity:   1000,
			Action:     ACTION_LOCK,
		},
		{
			Id:         "3",
			PlayerName: "TestAlt",
			TypeName:   "Oxygen",
			Quantity:   1500,
			Action:     ACTION_UNLOCK,
		},
	})

	if err != nil {
		t.Fatalf("Didn't expect error, but got %s", err)
	}

	// Make sure only EpicCyno has 500 left
	expectedContents := map[int][]InventoryStack{
		3683: {
			{PlayerName: "TestOther", Amount: 500},
		},
	}

	if !reflect.DeepEqual(expectedContents, h.inv.contents) {
		t.Fatalf("Expected %v, but got %v", expectedContents, h.inv.contents)
	}

	if len(debitMuts) != 1 {
		t.Fatalf("Expected %d ledger mutations, but got %d", 1, len(debitMuts))
	}

	if len(creditMuts) != 1 {
		t.Fatalf("Expected %d inventory mutations, but got %d", 1, len(creditMuts))
	}
}
