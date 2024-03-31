package services

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func (service AuthService) LogOut(authContext echo.Context) error {
	var err error

	request := new(models.SignOutUserRequest)
	authContext.Bind(request)

	token := helpers.GetTokenFromRequest(authContext)
	ctx := authContext.Request().Context()

	_, err = service.Rdb.Get(ctx, request.UserId).Result()
	var statusCode int = http.StatusNoContent
	if err != nil {
		if err == redis.Nil {
			service.Rdb.Set(ctx, request.UserId, token, 0)
		} else {
			statusCode = http.StatusInternalServerError
		}
		return authContext.JSON(statusCode, []string{err.Error()})
	}

	service.Rdb.Del(ctx, request.UserId)
	service.Rdb.Set(ctx, request.UserId, token, 0)
	return authContext.JSON(statusCode, nil)
}
