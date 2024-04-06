package services

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (service DomainService) GetDomains(domainContext echo.Context) error {
	var err error

	type domainDto struct {
		Name      string
		IsCustom  bool
		CreatedOn time.Time
	}

	request := new(models.GetDomainRequest)
	binder := &echo.DefaultBinder{}
	err = binder.BindHeaders(domainContext, request)
	if err != nil {
		return domainContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}
	err = binder.BindBody(domainContext, request)
	if err != nil {
		return domainContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	setDefaultDomainRequest(request)

	sortAndOrder := request.SortBy + " " + request.Order
	offset := (request.Page - 1) * request.PageLength
	query := fmt.Sprintf(`SELECT Name,IsCustom,CreatedOn FROM domains WHERE OrganizationId = '%s' AND IsDeprecated = %t ORDER BY %s LIMIT %d OFFSET %d`,
		request.OrganizationId, false, sortAndOrder, request.PageLength, offset)
	rows, err := service.Db.Query(query)
	if err != nil {
		return domainContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	defer rows.Close()

	var domains []domainDto
	for rows.Next() {
		var domain domainDto
		err = rows.Scan(&domain.Name, &domain.IsCustom, &domain.CreatedOn)
		if err != nil {
			return domainContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		domains = append(domains, domain)
	}

	countQuery := `SELECT COUNT(1) FROM domains WHERE OrganizationId =? AND IsDeprecated=?`
	row := service.Db.QueryRow(countQuery, request.OrganizationId, false)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return domainContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	datas := lo.Map(domains, func(domain domainDto, _ int) models.GetDomainData {
		return models.GetDomainData{
			Name:      domain.Name,
			IsCustom:  domain.IsCustom,
			CreatedOn: domain.CreatedOn,
		}
	})

	totalPage := int(math.Ceil(float64(count) / float64(request.PageLength)))
	resp := models.GetDomainResponse{
		Domains:    datas,
		Page:       request.Page,
		TotalPage:  totalPage,
		Totals:     count,
		PageLength: len(datas),
	}

	return domainContext.JSON(http.StatusOK, resp)
}

func setDefaultDomainRequest(request *models.GetDomainRequest) {
	if request.Page < 1 {
		request.Page = 1
	}

	if request.PageLength < 10 {
		request.PageLength = 10
	}

	if request.SortBy == "" {
		request.SortBy = "Name"
	}

	if request.Order == "" {
		request.Order = "asc"
	}
}
