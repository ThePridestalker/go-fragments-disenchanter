package cmd

import (
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

type DataResult struct {
	Token   string
	BaseURL string
}

func DataRetrieval() DataResult {

	cmd := exec.Command("wmic", "PROCESS", "WHERE", "name='LeagueClientUx.exe'", "GET", "commandline")

	leagueClientProcessData, err := cmd.Output()

	if err != nil {
		log.Println("Error al ejecutar el comando:", err)
		return DataResult{}
	}

	// Convert to utf-8
	cmdOutput := string(leagueClientProcessData)

	// get the correct part of the output through regex
	portPattern := regexp.MustCompile(`--app-port=([0-9]*)`)
	passwordPattern := regexp.MustCompile(`--remoting-auth-token=([\w\-_]*)`)

	portMatch := portPattern.FindStringSubmatch(cmdOutput)
	passwordMatch := passwordPattern.FindStringSubmatch(cmdOutput)

	// getting the port and password from the matches
	var port, password string
	if len(portMatch) >= 2 {
		port = portMatch[1]
	}
	if len(passwordMatch) >= 2 {
		password = passwordMatch[1]
	}

	host := "https://127.0.0.1"
	username := "riot"

	// create the baseUrl
	baseUrl := fmt.Sprintf("%s:%s", host, port)

	// create the token
	token := stringToBase64(fmt.Sprintf("%s:%s", username, password))

	result := DataResult{
		Token:   token,
		BaseURL: baseUrl,
	}

	return result
}

func stringToBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}
