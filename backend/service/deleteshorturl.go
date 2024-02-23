package services

import (
	"errors"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service UrlService) DeletehortUrl(userContext echo.Context) error {
	var err error
	request := new(models.DeleteUrlRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	if len(request.ShortUrl) == 0 {
		return userContext.JSON(http.StatusBadRequest, errors.New("link should not be empty"))
	}

	result, _ := service.Db.Exec("UPDATE shortlinks SET IsDeprecated =? WHERE Hash =?", false, request.ShortUrl)
	rows, err := result.RowsAffected()
	if err == nil && rows > 0 {
		return userContext.JSON(http.StatusNoContent, nil)
	}

	return userContext.JSON(http.StatusNotFound, err)
}
