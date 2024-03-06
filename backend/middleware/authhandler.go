package middlewares

import (
	"database/sql"
	"net/http"
	"time"

	jwtauth "github.com/adewoleadenigbagbe/url-shortner-service/helpers/auth"
	"github.com/labstack/echo/v4"
)

func (appMiddleware *AppMiddleware) AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		id, err := jwtauth.ValidateJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "Authentication required")
		}

		apikey := context.Request().Header["X-Api-Key"]
		if len(apikey) == 0 {
			return context.JSON(http.StatusBadRequest, "User ApiKey missing in the header")
		}

		db := context.Request().Context().Value(Db).(*sql.DB)
		if !isCurrentUser(db, id, apikey[0]) {
			return context.JSON(http.StatusUnauthorized, "invalid api key")
		}
		return next(context)
	}
}

func isCurrentUser(db *sql.DB, id, key string) bool {
	query := `SELECT COUNT(1) FROM users JOIN userkeys on users.Id = userkeys.UserId 
	WHERE users.Id =? AND userkeys.ApiKey =? AND userkeys.IsActive =? AND userkeys.ExpirationDate >?`
	var count int64
	err := db.QueryRow(query, id, key, true, time.Now()).Scan(&count)
	if err != nil {
		return false
	}
	return count == 1
}
