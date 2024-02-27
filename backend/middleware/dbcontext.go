package middlewares

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
)

type ContextType int

const (
	Db ContextType = 1
)

type AppMiddleware struct {
	Db *sql.DB
}

func (appMiddleware *AppMiddleware) SetDbContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), Db, appMiddleware.Db)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
