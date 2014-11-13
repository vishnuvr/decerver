package main

import (
	"github.com/eris-ltd/decerver/hooks/server"
	ghh "github.com/eris-ltd/decerver/hooks/server/ghhandler"
)

// This program will create a webserver and handler for github webhooks.
// ctrl-C in the terminal window to shut it down (gracefully).
func main() {
	ws := server.NewWebServer()
	// Pass the github webhook handler.
	ghHandler := ghh.NewHandler()
	ws.Start(ghHandler)
}
