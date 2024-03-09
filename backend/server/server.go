package server

import (
	"context"
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

func (server *ApplicationServer) Serve() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//set echo log
	server.BaseApp.Echo.Logger.SetLevel(log.INFO)

	//cors middleware
	server.BaseApp.Echo.Use(middleware.CORS())

	//Register Routes
	routes.RegisterRoutes(server.BaseApp, server.AppMiddleWare)

	// Start server
	go func() {
		if err := server.BaseApp.Echo.Start(":8653"); err != nil && err != http.ErrServerClosed {
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
