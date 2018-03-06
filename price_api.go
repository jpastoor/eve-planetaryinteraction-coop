package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type PriceAPI interface {
	FetchPrice(ty Type) (float32, error)
}

type EveMarketerAPI struct {
	client *http.Client
}

type EveMarketerRsp []EveMarketerRspType

type EveMarketerRspType struct {
	Buy  EveMarketerRspTypeStat `json:"buy"`
	Sell EveMarketerRspTypeStat `json:"sell"`
}

type EveMarketerRspTypeStat struct {
	FivePercent float32 `json:"fivePercent"`
}

/**
Currently fetches the Buy 5% price
 */
func (l *EveMarketerAPI) FetchPrice(ty Type) (float32, error) {
	// TODO Add Caching
	// TODO add grouping of typeIds in a single call

	url := fmt.Sprintf("https://api.evemarketer.com/ec/marketstat/json?usesystem=30000142&typeid=%d", ty.TypeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	rsp, err := l.client.Do(req)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}

	parsedRsp := EveMarketerRsp{}
	if err := json.Unmarshal(body, &parsedRsp); err != nil {
		return 0, err
	}

	return parsedRsp[0].Buy.FivePercent, nil
}
