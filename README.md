# GoRAT v1.1.0
Remote Access Tool written in Golang (client) and JavaScript using Node.js (server).

### Disclaimer (READ BEFORE CONTINUING)
This tool was made for educational/research purposes only. The author **IS NOT** responsible nor liable for how you use anything provided here. The author and/or anyone affiliated with the project **WILL NOT** be responsible nor liable for any losses or damages caused by this tool. By using GoRAT **YOU** and only you, have full responsibility for **ANYTHING** you cause. By using or downloading **GoRAT** you automatically **AGREE TO USE GoRAT AT YOUR OWN RISK**. **GoRAT** is **ONLY INTENDED** to be used on **YOUR OWN** pentesting machines/labs or with **EXPLICIT** consent from the owner of the property being pentested.
With this disclaimer, the author and/or anyone affiliated with **GoRAT** is **FULLY** exempt from anything you or somebody can cause with it.

### Features
* Retrieving public IP
* Retrieving machine name and username
* Executable persistency (using regedit)
* **NEW** Downloading files (client-side)
* **NEW** Downloading and opening files (client-side) 
* **NEW** Uploading files from client to server
* **NEW** Manually broadcast to every client by sending commands from web (default port: 9899)

### Installation
###### Server installation
1. Install [Node.js](https://nodejs.org/en/)
2. `cd gorat\server`
3. `npm i`
4. `node server.js` (use forever or nodemon)
###### Client installation
1. Install [Golang](https://golang.org/doc/install)
2. `cd gorat\client`
3. `go build app.go`
