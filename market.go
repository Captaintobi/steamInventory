package main

import (
	"encoding/json"
	"net/http"
)

type Price struct {
	Success     bool   `json:"success"`
	LowestPrice string `json:"lowest_price"`
	Volume      string `json:"volume"`
	MedianPrice string `json:"median_price"`
}

func getPrice(appid int64, marketname string) (*Price, error) {
	resp, err := http.Get("http://steamcommunity.com/market/priceoverview/?currency=1&appid=" + string(appid) + "&market_hash_name=" + marketname)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	var price Price
	json.NewDecoder(resp.Body).Decode(&price)
	return &price, nil
}
