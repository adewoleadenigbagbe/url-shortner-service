package main

import (
	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	"github.com/labstack/gommon/log"
)

func main() {
	app, err := core.ConfigureApp()
	if err != nil {
		log.Fatal(err)
	}

	server := ApplicationServer{
		App: app,
	}

	server.Serve()
}
