# GoRAT v1.1.0
Remote Access Tool written in Golang (client) and JavaScript (server).

### Disclaimer (READ BEFORE CONTINUING)
This tool was made for educational/research purposes only. The author **IS NOT** responsible nor liable for how you use anything provided here. The author and/or anyone affiliated with the project **WILL NOT** be responsible nor liable for any losses or damages caused by this tool. By using GoRAT **YOU** and only you, have full responsibility for **ANYTHING** you cause. By using or downloading **GoRAT** you automatically **AGREE TO USE GoRAT AT YOUR OWN RISK**. **GoRAT** is **ONLY INTENDED** to be used on **YOUR OWN** pentesting machines/labs or with **EXPLICIT** consent from the owner of the property being pentested.
With this disclaimer, the author and/or anyone affiliated with **GoRAT** is **FULLY** exempt from anything you or somebody can cause with it.

### Features
* Retrieving computer public IP
* Retrieving machine name and username
* **NEW** Executable persistency (using regedit)
* **NEW** Downloading files (client-side)
* **NEW** Downloading and opening files (client-side) 
* **NEW** Uploading files from client to server

### Installation
###### Server installation
1. Install PHP ([Ubuntu 18.04](https://linuxize.com/post/how-to-install-php-on-ubuntu-18-04/), [Windows 10 (XAMPP)](https://www.apachefriends.org/download.html))
2. Go to your server's htdocs and place [server.php](https://github.com/tinopai/gorat/blob/master/server/server.php) there
###### Client installation
1. Install [Golang](https://golang.org/doc/install)
2. `cd gorat\client`
3. `go build app.go`
