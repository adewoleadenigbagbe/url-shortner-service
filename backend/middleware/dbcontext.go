package middlewares

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type ContextType int

const (
	Db ContextType = 1
)

type AppMiddleware struct {
	Db  *sql.DB
	Rdb *redis.Client
}

// func (appMiddleware *AppMiddleware) SetDbContext(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		ctx := context.WithValue(c.Request().Context(), Db, appMiddleware.Db)
// 		c.SetRequest(c.Request().WithContext(ctx))
// 		return next(c)
// 	}
// }
