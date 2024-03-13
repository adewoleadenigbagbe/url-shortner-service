package main

import (
	"fmt"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	middlewares "github.com/adewoleadenigbagbe/url-shortner-service/middleware"
	"github.com/adewoleadenigbagbe/url-shortner-service/server"
	"github.com/labstack/gommon/log"
)

func main() {
	adminUUID := sequentialguid.NewSequentialGuid().String()
	userUUID := sequentialguid.NewSequentialGuid().String()

	fmt.Println(adminUUID)
	fmt.Println(userUUID)

	fmt.Println(time.Now())

	app, err := core.ConfigureAppDependencies()
	if err != nil {
		log.Fatal(err)
	}

	server := server.ApplicationServer{
		BaseApp: app,
		AppMiddleWare: &middlewares.AppMiddleware{
			Db:  app.Db,
			Rdb: app.Rdb,
		},
	}

	server.Serve()
}
