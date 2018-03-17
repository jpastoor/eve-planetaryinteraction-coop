package main

type Ledger struct {
	priceApi PriceAPI
}

type LedgerMutation struct {
	TypeId     int
	TypePrice  float32
	Change     float32
	PlayerName string
}

/**
This type will be persisted in the database at every commit
 */
type LedgerState struct {
	CommitId   int    `gorm:"primary_key auto_increment=false"`
	TypeId     int    `gorm:"primary_key auto_increment=false"`
	PlayerName string `gorm:"primary_key auto_increment=false"`
	TypePrice  float32
	Change     float32
}

func NewLedger(priceApi PriceAPI) (*Ledger) {
	return &Ledger{priceApi: priceApi}
}

func (l *Ledger) HandleMutations(debit []InventoryMutation, credit []InventoryMutation) ([]LedgerMutation, error) {
	var lml []LedgerMutation

	// Find all unique typeIds and their prices
	typeIds := make(map[int]struct{})
	for _, mut := range append(debit, credit...) {
		if _, exists := typeIds[mut.TypeId]; !exists {
			typeIds[mut.TypeId] = struct{}{}
		}
	}

	uniqueTypeIds := []int{}
	for key := range typeIds {
		uniqueTypeIds = append(uniqueTypeIds, key)
	}

	typeIdsWithPrices, err := l.priceApi.FetchPrices(uniqueTypeIds)
	if err != nil {
		return nil, err
	}

	for _, mut := range append(debit, credit...) {
		price := typeIdsWithPrices[mut.TypeId]
		lml = append(lml, LedgerMutation{
			TypePrice:  price,
			TypeId:     mut.TypeId,
			Change:     price * float32(mut.Change),
			PlayerName: mut.PlayerName,
		})
	}

	return lml, nil
}


func (l *Ledger) CalculateLedgerSummary(mutations []LedgerMutation) []GetLedgerRspItem {
	ledgerByPlayer := make(map[string]*GetLedgerRspItem)
	for _, mutation := range mutations {
		if _, exists := ledgerByPlayer[mutation.PlayerName]; !exists {
			ledgerByPlayer[mutation.PlayerName] = &GetLedgerRspItem{
				PlayerName: mutation.PlayerName,
				Amount:     float64(mutation.Change),
			}
		} else {
			item := ledgerByPlayer[mutation.PlayerName]
			item.Amount += float64(mutation.Change)
		}
	}

	ledgerSummary := []GetLedgerRspItem{}
	for _, ledgerRspItem := range ledgerByPlayer {
		ledgerSummary = append(ledgerSummary, *ledgerRspItem)
	}
	return ledgerSummary
}