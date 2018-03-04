package main

import (
	"testing"
	"net/http"
)

func TestEveMarketerAPI_FetchPrice(t *testing.T) {
	api := EveMarketerAPI{
		client: &http.Client{},
	}

	price, err := api.FetchPrice(Type{TypeId: 3683})
	if err != nil {
		t.Fatalf("Did not expect an error, but got %s", err)
	}

	if price > 1000 || price < 100 {
		t.Fatalf("Expected price between 100-1000, but price is %d", price)
	}
}
