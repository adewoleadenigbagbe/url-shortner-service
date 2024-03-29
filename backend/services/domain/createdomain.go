package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
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
	binder := &echo.DefaultBinder{}
	binder.BindHeaders(domainContext, request)
	binder.BindBody(domainContext, request)

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

	if request.IsCustom {
		planType, noDomains, err := GetCustomDomainCount(service.Db, request.OrganizationId)
		if err != nil {
			return domainContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}

		if planType.Valid {
			var respErr string
			if planType.Val == enums.Free && noDomains >= models.Free_Plan_Domain_Limit {
				respErr = fmt.Sprintf("you have exceeded the domain limit for this plan type : %d", models.Free_Plan_Domain_Limit)
				return domainContext.JSON(http.StatusBadRequest, []string{respErr})
			} else if planType.Val == enums.Team && noDomains >= models.Team_Plan_Domain_Limit {
				respErr = fmt.Sprintf("you have exceeded the domain limit for this plan type : %d", models.Team_Plan_Domain_Limit)
				return domainContext.JSON(http.StatusBadRequest, []string{respErr})
			}
		}
	}

	now := time.Now()
	domainId := sequentialguid.NewSequentialGuid().String()
	_, err = service.Db.Exec("INSERT INTO domains VALUES(?,?,?,?,?,?,?,?);",
		domainId, request.Name, request.IsCustom, request.OrganizationId, now, now, request.UserId, false)

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

func GetCustomDomainCount(db *sql.DB, id string) (helpers.Nullable[enums.PayPlan], int, error) {
	var count int
	var planType helpers.Nullable[enums.PayPlan] //

	query := `SELECT payplans.Type,COUNT(domains.Id) AS domaincount FROM payplans 
	JOIN organizationpayplans ON payplans.Id = organizationpayplans.PayPlanId
	LEFT JOIN domains ON organizationpayplans.OrganizationId = domains.OrganizationId
	WHERE organizationpayplans.OrganizationId =?
	AND domains.IsDeprecated=?
	AND domains.IsCustom=?
	AND payplans.IsLatest=?
	AND organizationpayplans.IsLatest=?`

	err := db.QueryRow(query, id, false, true, true, true).Scan(&planType, &count)
	if err != nil {
		return helpers.NewNullable(enums.PayPlan(0), false), 0, err
	}
	return planType, count, nil
}
