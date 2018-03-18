package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
)

type PriceAPI interface {
	FetchPrices(typeIds []int) (map[int]float32, error)
}

type EveMarketerAPI struct {
	client *http.Client
	cache  map[int]float32 // TODO Cache should have some notion of time. (i.e. keep prices for 6 hours)
}

type EveMarketerRsp []EveMarketerRspType

type EveMarketerRspType struct {
	Buy  EveMarketerRspTypeStat `json:"buy"`
	Sell EveMarketerRspTypeStat `json:"sell"`
}

type EveMarketerRspTypeStat struct {
	ForQuery    EveMarketerRspForQuery `json:"forQuery"`
	FivePercent float32                `json:"fivePercent"`
}

type EveMarketerRspForQuery struct {
	Types []int `json:"types"`
}

/**
Currently fetches the Buy 5% price
 */
func (l *EveMarketerAPI) FetchPrices(typeIds []int) (map[int]float32, error) {
	output := make(map[int]float32)

	var typeIdsStr []string
	for _, typeID := range typeIds {
		// If we have a cache
		if price, exists := l.cache[typeID]; exists {
			output[typeID] = price
		} else {
			// When there is no cache, add it to the lists of types to fetch
			typeIdsStr = append(typeIdsStr, strconv.Itoa(typeID))
		}
	}

	if len(typeIdsStr) == 0 {
		return  output, nil
	}

	url := fmt.Sprintf("https://api.evemarketer.com/ec/marketstat/json?usesystem=30000142&typeid=%s", strings.Join(typeIdsStr, ","))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return output, err
	}

	rsp, err := l.client.Do(req)
	if err != nil {
		return output, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error while fetching price API [%d] %s on URL %s", rsp.StatusCode, rsp.Status, url)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return output, err
	}

	parsedRsp := EveMarketerRsp{}
	if err := json.Unmarshal(body, &parsedRsp); err != nil {
		return output, err
	}

	for _, typeStat := range parsedRsp {
		typeId := typeStat.Buy.ForQuery.Types[0]
		output[typeId] = typeStat.Buy.FivePercent
		l.cache[typeId] = output[typeId]
	}

	return output, nil
}
