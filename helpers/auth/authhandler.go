package jwtauth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		id, err := ValidateJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "Authentication required")
		}

		apikey := context.Request().Header["X-Api-Key"][0]
		if len(apikey) == 0 {
			return context.JSON(http.StatusBadRequest, "user ApiKey Missing in the header")
		}

		if !isCurrentUser(id, apikey) {
			return context.JSON(http.StatusUnauthorized, "invalid api key")
		}
		return next(context)
	}
}

func isCurrentUser(id, key string) bool {
	return false
}
