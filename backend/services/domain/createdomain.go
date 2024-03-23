package services

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

const (
	DuplicateName = "Duplicate... Name already exist"
)

type DomainService struct {
	Db *sql.DB
}

func (service DomainService) CreateDomain(domainContext echo.Context) error {
	var (
		err error
	)
	request := new(models.CreateDomainRequest)
	err = domainContext.Bind(request)
	if err != nil {
		return domainContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	errs := validateDomainRequest(*request)
	if len(errs) > 0 {
		valErrors := lo.Map(errs, func(er error, index int) string {
			return er.Error()
		})
		return domainContext.JSON(http.StatusBadRequest, valErrors)
	}

	var count int64
	err = service.Db.QueryRow("SELECT COUNT(1) FROM domains WHERE Name =?", request.Name).Scan(&count)
	if err != nil {
		return domainContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if count > 0 {
		return domainContext.JSON(http.StatusBadRequest, []string{DuplicateName})
	}

	now := time.Now()
	domainId := sequentialguid.NewSequentialGuid().String()
	_, err = service.Db.Exec("INSERT INTO domains VALUES(?,?,?,?,?,?);",
		domainId, request.Name, now, now, false, request.UserId)

	if err != nil {
		return domainContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	return domainContext.JSON(http.StatusCreated, models.CreateDomainResponse{DomainId: domainId, Name: request.Name})
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
