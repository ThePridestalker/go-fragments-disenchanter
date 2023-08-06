package main

import (
	"fmt"
	"fragments-disenchanter/cmd"
)

func main() {
	result := cmd.DataRetrieval()

	fmt.Println("Token:", result.Token)
	fmt.Println("URL:", result.BaseURL)
}
