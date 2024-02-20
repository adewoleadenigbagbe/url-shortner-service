package main

import (
	"log"

	"github.com/adewoleadenigbagbe/url-shortner-service/core"
)

func main() {
	app, err := core.ConfigureApp()
	if err != nil {
		log.Fatal(err)
	}

	server := core.ApplicationServer{
		App: app,
	}

	server.Serve()
}
