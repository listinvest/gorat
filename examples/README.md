# GoRAT Examples (not finished)

### List of GoRAT possible RFS requests
#### RFS:DOFile
Download and open files from the internet.
###### Format
`RFS:DOFile[<&>]ip:<machine public ip/all>[<&>]uri:<file download link>[<&>]fn:<file name.extension>`
###### Example
`RFS:DOFile[<&>]ip:all[<&>]uri:https://raw.githubusercontent.com/tinopai/gorat/master/examples/test.bat[<&>]fn:test.bat`

#### RFS:Screenshot
Take a screenshot and send it to the server encoded in Base 64, then save it as `Screenshot{unix timestamp}{machine name}{random string}.png`
###### Format
`RFS:Screenshot[<&>]ip:<machine public ip/all>`
###### Example
`RFS:Screenshot[<&>]ip:all`

#### RFS:OFile
Open a file located in 

#### RFS:UFile
Tries to upload a file from client to server, returns file data if the action was successful and returns 'Couldnt upload file: {error}[]'
