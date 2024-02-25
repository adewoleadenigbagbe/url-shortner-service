package services

import (
	"errors"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service UrlService) DeleteShortUrl(urlContext echo.Context) error {
	var err error
	request := new(models.DeleteUrlRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, err.Error())
	}

	if len(request.ShortUrl) == 0 {
		return urlContext.JSON(http.StatusBadRequest, errors.New("link should not be empty"))
	}

	result, _ := service.Db.Exec("UPDATE shortlinks SET IsDeprecated =? WHERE Hash =?", false, request.ShortUrl)
	rows, err := result.RowsAffected()
	if err == nil && rows > 0 {
		return urlContext.JSON(http.StatusNoContent, nil)
	}

	return urlContext.JSON(http.StatusNotFound, err)
}
