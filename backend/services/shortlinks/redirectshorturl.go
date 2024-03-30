package services

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service UrlService) RedirectShort(urlContext echo.Context) error {
	var err error
	request := new(models.RedirectShortRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	if len(request.ShortUrl) == 0 {
		return urlContext.JSON(http.StatusBadRequest, []string{"shortlink is invalid"})
	}

	var originalUrl string
	query := `SELECT shortlinks.OriginalUrl
				FROM shortlinks JOIN domains ON shortlinks.DomainId = domains.Id 
				WHERE shortlinks.Hash=? AND shortlinks.IsDeprecated=? AND domains.IsDeprecated=?`
	row := service.Db.QueryRow(query, request.ShortUrl, false, false)
	if err = row.Scan(&originalUrl); errors.Is(err, sql.ErrNoRows) {
		return urlContext.JSON(http.StatusNotFound, []string{err.Error()})
	}

	id := sequentialguid.NewSequentialGuid().String()
	createdOn := time.Now()
	_, err = service.Db.Exec("INSERT INTO accesslogs VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);",
		id, request.ShortUrl, request.Country, request.TimeZone, request.City,
		request.Os, request.Browser, request.UserAgent, request.Platform,
		request.IpAddress, request.Method, request.Status, request.OrganizationId, createdOn, false)

	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	urlContext.Response().Header().Set("Location", originalUrl)
	return urlContext.JSON(http.StatusFound, nil)
}
