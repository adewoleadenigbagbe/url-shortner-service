package services

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		return urlContext.JSON(http.StatusNotFound, []string{"link does not exist"})
	}

	setFieldToTitleCase(request)

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

func setFieldToTitleCase(request *models.RedirectShortRequest) {
	if request.Browser.Valid {
		request.Browser.Val = cases.Title(language.English, cases.Compact).String(request.Browser.Val)
	}

	if request.City.Valid {
		request.City.Val = cases.Title(language.English, cases.Compact).String(request.City.Val)
	}

	if request.Country.Valid {
		request.Country.Val = cases.Title(language.English, cases.Compact).String(request.Country.Val)
	}

	if request.Os.Valid {
		request.Os.Val = cases.Title(language.English, cases.Compact).String(request.Os.Val)
	}

	if request.Platform.Valid {
		request.Platform.Val = cases.Title(language.English, cases.Compact).String(request.Platform.Val)
	}
}
