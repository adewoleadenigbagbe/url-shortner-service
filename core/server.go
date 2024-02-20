package core

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/gommon/log"
)

type ApplicationServer struct {
	App *BaseApp
}

func (server *ApplicationServer) Serve() {
	server.App.echo.Logger.SetLevel(log.INFO)
	RegisterRoutes(server.App)

	// Start server
	go func() {
		if err := server.App.echo.Start(":8653"); err != nil && err != http.ErrServerClosed {
			server.App.echo.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.App.echo.Shutdown(ctx); err != nil {
		server.App.echo.Logger.Fatal(err)
	}
}
