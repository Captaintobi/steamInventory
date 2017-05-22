package main

import (
	newApi "github.com/Philipp15b/go-steamapi"
)

type MyInventory struct {
	Result Inventory
}
type Inventory struct {
	Assets []struct {
		Appid      string `json:"appid"`
		Contextid  string `json:"contextid"`
		Assetid    string `json:"assetid"`
		Classid    string `json:"classid"`
		Instanceid string `json:"instanceid"`
		Amount     string `json:"amount"`
	} `json:"assets"`
	Descriptions []struct {
		Appid           int    `json:"appid"`
		Classid         string `json:"classid"`
		Instanceid      string `json:"instanceid"`
		Currency        int    `json:"currency"`
		BackgroundColor string `json:"background_color"`
		IconURL         string `json:"icon_url"`
		IconURLLarge    string `json:"icon_url_large,omitempty"`
		Descriptions    []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
			Color string `json:"color,omitempty"`
		} `json:"descriptions"`
		Tradable int `json:"tradable"`
		Actions  []struct {
			Link string `json:"link"`
			Name string `json:"name"`
		} `json:"actions,omitempty"`
		Name           string `json:"name"`
		NameColor      string `json:"name_color"`
		Type           string `json:"type"`
		MarketName     string `json:"market_name"`
		MarketHashName string `json:"market_hash_name"`
		MarketActions  []struct {
			Link string `json:"link"`
			Name string `json:"name"`
		} `json:"market_actions,omitempty"`
		Commodity                 int `json:"commodity"`
		MarketTradableRestriction int `json:"market_tradable_restriction"`
		Marketable                int `json:"marketable"`
		Tags                      []struct {
			Category              string `json:"category"`
			InternalName          string `json:"internal_name"`
			LocalizedCategoryName string `json:"localized_category_name"`
			LocalizedTagName      string `json:"localized_tag_name"`
			Color                 string `json:"color,omitempty"`
		} `json:"tags"`
	} `json:"descriptions"`
	TotalInventoryCount int `json:"total_inventory_count"`
	Success             int `json:"success"`
	Rwgrsn              int `json:"rwgrsn"`
}

func friendToMap(id uint64, filter newApi.Relationship, apiKey string) (map[string]uint64, error) {
	friendList, err := newApi.GetFriendsList(id, filter, apiKey)
	if err != nil {
		return nil, err
	}
	myFriends := make(map[string]uint64)
	var myFriendIds []uint64
	//making a slice of steam ids so that GetPlayersummaries can use it
	for _, friend := range friendList {

		myFriendIds = append(myFriendIds, friend.SteamID)

	}
	for i := 1; i < (len(myFriendIds)/100)+1; i++ {

		friendSummaries, err := newApi.GetPlayerSummaries(myFriendIds[100*(i-1):100*i], apiKey)
		if err != nil {
			return nil, err
		}

		for _, friendSummary := range friendSummaries {
			//fmt.Println(friendSummary.PersonaName, friendSummary.SteamID)
			myFriends[friendSummary.PersonaName] = friendSummary.SteamID

		}

	} //missng like 47 have to get the rest some how.
	friendSummaries, err := newApi.GetPlayerSummaries(myFriendIds[len(myFriends):], apiKey)
	if err != nil {
		return nil, err
	}

	for _, friendSummary := range friendSummaries {
		//fmt.Println(friendSummary.PersonaName, friendSummary.SteamID)
		myFriends[friendSummary.PersonaName] = friendSummary.SteamID

	}
	return myFriends, nil

}
