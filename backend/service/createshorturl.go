package services

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

const (
	expirySpan = 1
)

func (service UrlService) CreateShortUrl(urlContext echo.Context) error {
	var err error
	request := new(models.CreateUrlRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, err.Error())
	}

	errs := validateUrlRequest(*request)
	if len(errs) > 0 {
		return urlContext.JSON(http.StatusBadRequest, errs)
	}

	short := helpers.GenerateShortUrl(request.OriginalUrl)
	now := time.Now()
	expirationDate := now.AddDate(expirySpan, 0, 0)
	_, err = service.Db.Exec("INSERT INTO shortlinks VALUES(?,?,?,?,?,?,?,?,?,?);",
		short, request.OriginalUrl, request.DomainName, request.CustomAlias, sql.NullInt64{Valid: false}, now, now, expirationDate, false, request.UserId)

	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, err)
	}
	return urlContext.JSON(http.StatusCreated, models.CreateUrlResponse{ShortUrl: short})
}

func validateUrlRequest(request models.CreateUrlRequest) []error {
	var validationErrors []error

	if request.OriginalUrl == "" {
		validationErrors = append(validationErrors, errors.New("url is required"))
	}

	if request.UserId == "" {
		validationErrors = append(validationErrors, errors.New("userId is required"))
	}

	isValid := helpers.IsValidUrl(request.OriginalUrl)
	if !isValid {
		validationErrors = append(validationErrors, errors.New("invalid url"))
	}

	return validationErrors
}
