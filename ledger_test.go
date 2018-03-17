package main

import (
	"testing"
	"reflect"
)

func TestCalculateLedgerSummary(t *testing.T) {

	ledger := NewLedger(nil)

	muts := []LedgerMutation{
		{1, 12.3, 123, "Swaffeltje"},
		{1, 12.3, -23, "Swaffeltje"},
		{1, 12.3, -200, "Gebbetje"},
	}

	expected := []GetLedgerRspItem{
		{PlayerName: "Swaffeltje", Amount: 100},
		{PlayerName: "Gebbetje", Amount: -200},
	}
	summary := ledger.CalculateLedgerSummary(muts)

	if !reflect.DeepEqual(expected, summary) {
		t.Fatalf("Expected %v, but got %v", expected, summary)
	}
}
