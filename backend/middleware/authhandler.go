package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	jwtauth "github.com/adewoleadenigbagbe/url-shortner-service/helpers/auth"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func (appMiddleware *AppMiddleware) AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		id, err := jwtauth.ValidateUserJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "Authentication required")
		}

		tokenExist := checkForBlackListedTokens(context, appMiddleware.Rdb, id)
		if tokenExist {
			return context.JSON(http.StatusBadRequest, "Invalid Authourization Token")
		}

		apikey := context.Request().Header.Get("X-Api-Key")
		if len(apikey) == 0 {
			return context.JSON(http.StatusBadRequest, "User ApiKey missing in the header")
		}

		if !isCurrentUser(appMiddleware.Db, id, apikey) {
			return context.JSON(http.StatusUnauthorized, "invalid api key")
		}
		return next(context)
	}
}

func (appMiddleware *AppMiddleware) AuthorizeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		id, err := jwtauth.ValidateAdminRoleJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "You are not allowed to access this resource")
		}

		tokenExist := checkForBlackListedTokens(context, appMiddleware.Rdb, id)
		if tokenExist {
			return context.JSON(http.StatusBadRequest, "Invalid Authourization Token")
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

func checkForBlackListedTokens(context echo.Context, redisClient *redis.Client, id string) bool {
	token := jwtauth.GetTokenFromRequest(context)
	res, err := redisClient.Get(context.Request().Context(), id).Result()
	if err != redis.Nil {
		fmt.Println(err)
	}
	return res == token
}
