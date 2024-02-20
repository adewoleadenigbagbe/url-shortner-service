package core

import (
	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	services "github.com/adewoleadenigbagbe/url-shortner-service/service"

	"github.com/labstack/echo/v4"
)

type BaseApp struct {
	echo        *echo.Echo
	AuthService services.AuthService
}

func ConfigureApp() (*BaseApp, error) {
	db, err := database.ConnectToSQLite()
	if err != nil {
		return nil, err
	}
	app := &BaseApp{
		echo: echo.New(),
		AuthService: services.AuthService{
			Db: db,
		},
	}

	return app, nil
}
