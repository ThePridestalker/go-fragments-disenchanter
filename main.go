package main

import (
	"fmt"
	"fragments-disenchanter/api"
	"fragments-disenchanter/utils"
)

func main() {
	result := utils.DataRetrieval()

	if result.Token == "" || result.BaseURL == "" {
		fmt.Println("Token or URL is empty. Make sure your League of Legends client is up and running!")
		return
	}

	api.GetLootAndDisenchant(result.BaseURL, result.Token)
}
