package core

import (
	"database/sql"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	services "github.com/adewoleadenigbagbe/url-shortner-service/service"
	"github.com/labstack/echo/v4"
)

type BaseApp struct {
	Echo        *echo.Echo
	Db          *sql.DB
	AuthService services.AuthService
	UrlService  services.UrlService
}

func ConfigureAppDependencies() (*BaseApp, error) {
	db, err := database.ConnectToSQLite()
	if err != nil {
		return nil, err
	}
	app := &BaseApp{
		Echo: echo.New(),
		Db:   db,
		AuthService: services.AuthService{
			Db: db,
		},
		UrlService: services.UrlService{
			Db: db,
		},
	}

	return app, nil
}
