package middlewares

import (
	"database/sql"
	"net/http"

	jwtauth "github.com/adewoleadenigbagbe/url-shortner-service/helpers/auth"
	"github.com/labstack/echo/v4"
)

func AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		id, err := jwtauth.ValidateJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "Authentication required")
		}

		apikey := context.Request().Header["X-Api-Key"][0]
		if len(apikey) == 0 {
			return context.JSON(http.StatusBadRequest, "user ApiKey Missing in the header")
		}

		db := context.Request().Context().Value("db").(*sql.DB)
		if !isCurrentUser(db, id, apikey) {
			return context.JSON(http.StatusUnauthorized, "invalid api key")
		}
		return next(context)
	}
}

func isCurrentUser(db *sql.DB, id, key string) bool {
	return false
}
