package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	"github.com/adewoleadenigbagbe/url-shortner-service/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type ApplicationServer struct {
	App *core.BaseApp
}

func (server *ApplicationServer) Serve() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server.App.Echo.Logger.SetLevel(log.INFO)

	routes.RegisterRoutes(server.App)

	// Start server
	go func() {
		if err := server.App.Echo.Start(":8653"); err != nil && err != http.ErrServerClosed {
			server.App.Echo.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.App.Echo.Shutdown(ctx); err != nil {
		server.App.Echo.Logger.Fatal(err)
	}
}
