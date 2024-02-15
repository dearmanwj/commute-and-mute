# Commute and Mute

A lambda function to set activity status to commute and mute between home and work

## Building for lambda

Run the following to build with linux architecture with expected name:

`GOOS=linux GOARCH=amd64 go build -o bootstrap main.go`

Zip the `bootstrap` file and upload to aws


## TODO

Setup secret manager
Figure out who to host ui
Add delete/re login functionality
Add some build scripts?
