package middlewares

import (
	"context"

	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	"github.com/labstack/echo/v4"
)

type ContextType int

const (
	Db ContextType = 1
)

type AppMiddleware struct {
	app *core.BaseApp
}

func (appMiddleware *AppMiddleware) SetDbContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), Db, appMiddleware.app.Db)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
