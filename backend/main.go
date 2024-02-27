package main

import (
	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	middlewares "github.com/adewoleadenigbagbe/url-shortner-service/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	app, err := core.ConfigureAppDependencies()
	if err != nil {
		log.Fatal(err)
	}

	server := ApplicationServer{
		BaseApp: app,
		AppMiddleWare: &middlewares.AppMiddleware{
			Db: app.Db,
		},
	}

	server.Serve()
}
