package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

const (
	expirySpan = 1
)

func (service UrlService) CreateShortUrl(userContext echo.Context) error {
	var err error
	request := new(models.CreateUrlRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	errs := validateUrlRequest(*request)
	if len(errs) > 0 {
		return userContext.JSON(http.StatusBadRequest, errs)
	}

	var hashUrl string
	now := time.Now()
	expirationDate := now.AddDate(expirySpan, 0, 0)
	_, err = service.Db.Exec("INSERT INTO users VALUES(?,?,?,?,?,?);", hashUrl, request.OriginalUrl, now, now, expirationDate, request.UserId)
	if err != nil {
		return userContext.JSON(http.StatusInternalServerError, err)
	}
	return userContext.JSON(http.StatusOK, models.CreateUrlResponse{ShortUrl: hashUrl})
}

func validateUrlRequest(request models.CreateUrlRequest) []error {
	var validationErrors []error

	if request.OriginalUrl == "" {
		validationErrors = append(validationErrors, errors.New("url is required"))
	}

	if request.UserId == "" {
		validationErrors = append(validationErrors, errors.New("userId is required"))
	}

	return validationErrors
}
