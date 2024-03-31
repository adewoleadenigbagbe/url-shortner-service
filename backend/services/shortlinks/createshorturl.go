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
	expirySpan   = 1
	DuplicateUrl = "Duplicate... Url already exist"
)

type UrlService struct {
	Db *sql.DB
}

func (service UrlService) CreateShortLink(urlContext echo.Context) error {
	var err error
	request := new(models.CreateUrlRequest)
	binder := &echo.DefaultBinder{}
	binder.BindHeaders(urlContext, request)
	binder.BindBody(urlContext, request)

	errs := validateUrlRequest(*request)
	if len(errs) > 0 {
		valErrors := lo.Map(errs, func(er error, index int) string {
			return er.Error()
		})
		return urlContext.JSON(http.StatusBadRequest, valErrors)
	}

	var count int64
	query := `
		SELECT COUNT(1) FROM shortlinks
		JOIN domains ON shortlinks.DomainId = domains.Id
		WHERE shortlinks.OriginalUrl =? AND shortlinks.OrganizationId =?
		AND shortlinks.IsDeprecated =? AND domains.IsDeprecated =?
	`
	err = service.Db.QueryRow(query, request.OriginalUrl, request.OrganizationId, false, false).Scan(&count)
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if count > 0 {
		return urlContext.JSON(http.StatusBadRequest, []string{DuplicateUrl})
	}

	planType, linkCount, err := GetLinkCount(service.Db, request.OrganizationId)
	fmt.Println("plantype :", planType)
	fmt.Println("linkCount :", linkCount)
	fmt.Println("err :", err)

	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if planType.Valid && linkCount >= models.Free_Plan_Link_Limit {
		respErr := fmt.Sprintf("you have exceeded the link limit for this plan type : %d", models.Free_Plan_Link_Limit)
		return urlContext.JSON(http.StatusBadRequest, []string{respErr})
	}

	shortId := sequentialguid.NewSequentialGuid().String()
	short := helpers.GenerateShortLink(request.OriginalUrl)
	now := time.Now()
	expirationDate := now.AddDate(expirySpan, 0, 0)
	_, err = service.Db.Exec("INSERT INTO shortlinks VALUES(?,?,?,?,?,?,?,?,?,?,?);",
		shortId, short, request.OriginalUrl, request.DomainId, request.CustomAlias,
		now, now, expirationDate, request.OrganizationId, request.UserId, false)

	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	return urlContext.JSON(http.StatusCreated, models.CreateUrlResponse{Id: shortId, ShortUrl: short, DomainId: request.DomainId})
}

func validateUrlRequest(request models.CreateUrlRequest) []error {
	var validationErrors []error
	if request.DomainId == "" {
		validationErrors = append(validationErrors, errors.New("domainId is required"))
	}

	if request.UserId == "" {
		validationErrors = append(validationErrors, errors.New("userId is required"))
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

func GetLinkCount(db *sql.DB, id string) (helpers.Nullable[enums.PayPlan], int, error) {
	var count int
	var planType helpers.Nullable[enums.PayPlan] //

	query := `SELECT payplans.Type,COUNT(shortlinks.Id) AS linkcount FROM payplans 
	JOIN organizationpayplans ON payplans.Id = organizationpayplans.PayPlanId
	LEFT JOIN shortlinks ON organizationpayplans.OrganizationId = shortlinks.OrganizationId
	WHERE organizationpayplans.OrganizationId =?
	AND shortlinks.IsDeprecated=?
	AND payplans.IsLatest=?
	AND organizationpayplans.IsLatest=?`

	err := db.QueryRow(query, id, false, true, true).Scan(&planType, &count)
	if err != nil {
		return helpers.NewNullable(enums.PayPlan(0), false), 0, err
	}
	return planType, count, nil
}
