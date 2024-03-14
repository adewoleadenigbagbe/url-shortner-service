package services

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service UrlService) RedirectShort(urlContext echo.Context) error {
	var err error
	request := new(models.RedirectShortRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, err.Error())
	}

	if len(request.ShortUrl) == 0 {
		return urlContext.JSON(http.StatusBadRequest, errors.New("shortlink is invalid"))
	}

	var originalUrl string
	var hits int64
	var id string
	query := `SELECT shortlinks.Id,shortlinks.OriginalUrl, shortlinks.Hits FROM shortlinks JOIN domains ON shortlinks.DomainId = domains.Id 
	 WHERE shortlinks.Hash=? AND shortlinks.IsDeprecated=? AND domains.IsDeprecated=?`
	row := service.Db.QueryRow(query, request.ShortUrl, false, false)
	if err = row.Scan(&id, &originalUrl, &hits); errors.Is(err, sql.ErrNoRows) {
		return urlContext.JSON(http.StatusNotFound, err.Error())
	}

	//access the logs information
	tx, err := service.Db.Begin()
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err = tx.Exec("UPDATE shortlinks SET Hits =? WHERE Id =?", hits+1, id)
	if err != nil {
		tx.Rollback()
		return urlContext.JSON(http.StatusInternalServerError, err.Error())
	}

	//create accesslogs
	tx.Commit()

	urlContext.Response().Header().Set("Location", originalUrl)
	return urlContext.JSON(http.StatusFound, nil)
}
