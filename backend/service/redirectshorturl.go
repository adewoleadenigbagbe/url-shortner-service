package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service UrlService) LoginUser(urlContext echo.Context) error {
	var err error
	request := new(models.GetShortRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, err.Error())
	}

	if len(request.ShortUrl) == 0 {
		return urlContext.JSON(http.StatusBadRequest, errors.New("short links is invalid"))
	}

	var short string
	query := "SELECT OriginalUrl FROM shortlinks WHERE Hash=? AND IsDeprecated=?"
	row := service.Db.QueryRow(query, request.ShortUrl, false)
	if err = row.Scan(&short); errors.Is(err, sql.ErrNoRows) {
		return urlContext.JSON(http.StatusNotFound, err.Error())
	}

	urlContext.Response().Header().Set("Location", short)
	return urlContext.JSON(http.StatusFound, nil)
}
