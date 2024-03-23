package services

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

const (
	expirySpan   = 1
	DuplicateUrl = "Duplicate... Url already exist"
)

type UrlService struct {
	Db *sql.DB
}

func (service UrlService) CreateShortUrl(urlContext echo.Context) error {
	var err error
	request := new(models.CreateUrlRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, err.Error())
	}

	errs := validateUrlRequest(*request)
	if len(errs) > 0 {
		valErrors := lo.Map(errs, func(er error, index int) string {
			return er.Error()
		})
		return urlContext.JSON(http.StatusBadRequest, valErrors)
	}

	var count int64
	err = service.Db.QueryRow("SELECT COUNT(1) FROM shortlinks WHERE OriginalUrl =?", request.OriginalUrl).Scan(&count)
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, err.Error())
	}

	if count > 0 {
		return urlContext.JSON(http.StatusBadRequest, DuplicateUrl)
	}

	short := helpers.GenerateShortLink(request.OriginalUrl)
	now := time.Now()
	expirationDate := now.AddDate(expirySpan, 0, 0)
	_, err = service.Db.Exec("INSERT INTO shortlinks VALUES(?,?,?,?,?,?,?,?,?,?);",
		short, request.OriginalUrl, request.DomainId, request.CustomAlias, sql.NullInt64{Valid: false}, now, now, expirationDate, false, request.UserId)

	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, err)
	}
	return urlContext.JSON(http.StatusCreated, models.CreateUrlResponse{ShortUrl: short, DomainId: request.DomainId})
}

func validateUrlRequest(request models.CreateUrlRequest) []error {
	var validationErrors []error
	if request.UserId == "" {
		validationErrors = append(validationErrors, errors.New("userId is required"))
	}

	if request.DomainId == "" {
		validationErrors = append(validationErrors, errors.New("domainId is required"))
	}

	if request.OriginalUrl == "" {
		validationErrors = append(validationErrors, errors.New("url is required"))
	} else {
		isValid := helpers.IsValidUrl(request.OriginalUrl)
		if !isValid {
			validationErrors = append(validationErrors, errors.New("invalid url"))
		}
	}

	return validationErrors
}
