package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	middlewares "github.com/adewoleadenigbagbe/url-shortner-service/middleware"
	"github.com/adewoleadenigbagbe/url-shortner-service/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type ApplicationServer struct {
	BaseApp       *core.BaseApp
	AppMiddleWare *middlewares.AppMiddleware
}

func (server *ApplicationServer) start() {
	//set echo log
	server.BaseApp.Echo.Logger.SetLevel(log.INFO)

	//cors middleware
	server.BaseApp.Echo.Use(middleware.CORS())

	//Register Routes
	routes.RegisterRoutes(server.BaseApp, server.AppMiddleWare)

	// Start server
	go func() {
		port := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
		if err := server.BaseApp.Echo.Start(port); err != nil && err != http.ErrServerClosed {
			server.BaseApp.Echo.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.BaseApp.Echo.Shutdown(ctx); err != nil {
		server.BaseApp.Echo.Logger.Fatal(err)
	}
}

func InitializeAPI() {
	err := godotenv.Load(".env.example")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app, err := core.ConfigureAppDependencies()
	if err != nil {
		log.Fatal(err)
	}

	server := ApplicationServer{
		BaseApp: app,
		AppMiddleWare: &middlewares.AppMiddleware{
			Db:  app.Db,
			Rdb: app.Rdb,
		},
	}

	server.start()
}
