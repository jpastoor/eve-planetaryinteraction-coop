package main

import (
	"strings"
	"time"
	"fmt"
	"strconv"
)

type TransactionParser struct {
	typeFetcher TypeFetcher
}

func NewTransactionParser(typeFetcher TypeFetcher) (TransactionParser) {
	return TransactionParser{
		typeFetcher: typeFetcher,
	}
}

func (tp *TransactionParser) Parse(input string) (ts []Transaction, errs []error) {

	lines := strings.Split(input, "\n")
	for i, line := range lines {
		tabs := strings.Split(strings.TrimSpace(line), "\t")

		if len(tabs) != 8 {
			errs = append(errs, fmt.Errorf("expected 8 tabs on line %d, but got %d", i, len(tabs)))
			continue
		}

		// Parse Time
		creationDate, err := time.Parse("2006.01.02 15:04:05", tabs[0])
		if err != nil {
			errs = append(errs, fmt.Errorf("could not parse time on line %d: %s", i, err))
			continue
		}

		action := tabs[4]
		status := tabs[5]
		typeName := tabs[6]
		
		// Skip some specific cases
		if action == "Assembled" || action == "Set Name" || action == "Configure" || typeName == "Station Container" || typeName == "Station Vault Container" {
			continue
		}

		// Parse amount
		amount, err := strconv.Atoi(tabs[7])
		if err != nil {
			errs = append(errs, fmt.Errorf("could not parse amount on line %d: %s", i, err))
			continue
		}

		ty, err := tp.typeFetcher.getTypeByName(typeName)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not parse type on line %d: %s", i, err))
			continue
		}

		tp := Transaction{
			CreationDate: creationDate,
			Location:     tabs[1],
			SubLocation:  tabs[2],
			PlayerName:   tabs[3],
			Action:       action,
			Status:       status,
			TypeName:     ty.TypeName,
			Quantity:     amount,
		}

		tp.Id = tp.hash()

		ts = append(ts, tp)
	}

	return ts, errs
}
