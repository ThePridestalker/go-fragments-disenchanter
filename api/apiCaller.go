package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
	endpointLOOT       = "/lol-loot/v1/player-loot"
	endpointDISENCHANT = "/lol-loot/v1/recipes/CHAMPION_RENTAL_disenchant/craft?repeat="
)

type Shard struct {
	DisenchantLootName string `json:"disenchantLootName"`
	ItemStatus         string `json:"itemStatus"`
	Count              uint   `json:"count"`
	LootName           string `json:"lootName"`
}

var customClient *http.Client

func init() {
	// Create a custom HTTP client with disabled SSL verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	customClient = &http.Client{Transport: transport}
}

func getShards(url, token string) ([]Shard, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+token)

	resp, err := customClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var loot []Shard

	if err := json.NewDecoder(resp.Body).Decode(&loot); err != nil {
		return nil, err
	}

	return loot, nil

}

func sendDisenchantReq(url, champShard, token string) (string, error) {

	// create a single element array cause we cant send more at once
	finalShard := [1]string{champShard}
	data, err := json.Marshal(finalShard)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+token)

	resp, err := customClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var finalResponse string
	if err := json.NewDecoder(resp.Body).Decode(&finalResponse); err != nil {
		return "", err
	}

	return finalResponse, nil

}

func GetLootAndDisenchant(baseURL, token string) {
	urlGET := baseURL + endpointLOOT
	urlPOST := baseURL + endpointDISENCHANT

	loot, err := getShards(urlGET, token)
	if err != nil {
		log.Println("Error getting loot:", err)
		return
	}

	for _, shard := range loot {
		if shard.DisenchantLootName == "CURRENCY_champion" && shard.ItemStatus == "OWNED" {
			_, err := sendDisenchantReq(urlPOST+strconv.FormatUint(uint64(shard.Count), 10), shard.LootName, token)
			if err != nil {
				log.Println("Error disenchanting:", err)
				// Continue with the next item
				continue
			}
		}
	}
}
