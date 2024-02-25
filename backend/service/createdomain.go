package services

import (
	"errors"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service DomainService) CreateDomain(domainContext echo.Context) error {
	var err error
	request := new(models.CreateDomainRequest)
	err = domainContext.Bind(request)
	if err != nil {
		return domainContext.JSON(http.StatusBadRequest, err.Error())
	}

	errs := validateDomainRequest(*request)
	if len(errs) > 0 {
		return domainContext.JSON(http.StatusBadRequest, errs)
	}

	now := time.Now()
	domainId := sequentialguid.NewSequentialGuid().String()
	_, err = service.Db.Exec("INSERT INTO domains VALUES(?,?,?,?,?,?);",
		domainId, request.Name, now, now, false, request.UserId)

	if err != nil {
		return domainContext.JSON(http.StatusInternalServerError, err)
	}
	return domainContext.JSON(http.StatusCreated, models.CreateUrlResponse{ShortUrl: domainId})
}

func validateDomainRequest(request models.CreateDomainRequest) []error {
	var validationErrors []error

	if request.Name == "" {
		validationErrors = append(validationErrors, errors.New("name is required"))
	}

	if request.UserId == "" {
		validationErrors = append(validationErrors, errors.New("userId is required"))
	}
	return validationErrors
}
