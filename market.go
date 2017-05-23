package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Marketprice struct {
	Success     bool   `json:"success"`
	LowestPrice string `json:"lowest_price"`
	Volume      string `json:"volume"`
	MedianPrice string `json:"median_price"`
}

func getPrice(appid int, marketname string) (*Marketprice, error) {
	price := new(Marketprice)
	base, err := url.Parse("http://steamcommunity.com/market/priceoverview/?")
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	s := strconv.Itoa(appid)

	v.Set("appid", s)
	v.Add("currency", string(1))
	v.Add("market_hash_name", marketname)
	base.RawQuery = v.Encode()
	//resp, err := http.Get("http://steamcommunity.com/market/priceoverview/?appid=730&currency=1&market_hash_name=StatTrak%E2%84%A2%20M4A1-S%20|%20Hyper%20Beast%20(Minimal%20Wear)")
	fmt.Println("we use", base.String())
	resp, err := http.Get(base.String())
	if err != nil {
		return nil, err

	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(price)
	return price, nil
}
