package main

import (
	"testing"
	"net/http"
)

func TestEveMarketerAPI_FetchPrice(t *testing.T) {
	api := EveMarketerAPI{
		client: &http.Client{},
	}

	prices, err := api.FetchPrices([]int{2393, 3683})
	if err != nil {
		t.Fatalf("Did not expect an error, but got %s", err)
	}

	if prices[2393] > 1000 || prices[2393] < 100 {
		t.Fatalf("Expected price between 100-1000, but price is %f", prices[2393])
	}

	if prices[3683] > 1000 || prices[3683] < 100 {
		t.Fatalf("Expected price between 100-1000, but price is %f", prices[3683])
	}
}
