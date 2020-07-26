# GoRAT
Remote Access Tool written in Golang (client) and PHP (server).

### Disclaimer
This tool was made for educational purposes only.
I am not responsible nor condone other people using this tool for malicious purposes.
This program should be only ran on computers you own or have permissions to pentest.

### Features
* Retrieving computer public IP
* Retrieving machine name and username
* Sending that data to Discord via webhooks

### Installation
**Server installation** | Changing to JavaScript most-likely
1. Install PHP ([Ubuntu 18.04](https://linuxize.com/post/how-to-install-php-on-ubuntu-18-04/), [Windows 10 (XAMPP)](https://www.apachefriends.org/download.html))
2. Go to your server's htdocs and place [server.php](https://github.com/tinopai/gorat/blob/master/server/server.php) there
**Client installation**
1. Install [Golang](https://golang.org/doc/install)
2. `cd` to where GoRAT is located
3. Build the application by running `go build app.go`
