package services

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers/sqltype"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (urlService UrlService) GetShortLinks(urlContext echo.Context) error {
	var (
		err error
	)
	type shortDto struct {
		Hash           string
		OriginalUrl    string
		DomainName     string
		Alias          sqltype.Nullable[string]
		CreatedOn      time.Time
		ExpirationDate time.Time
	}

	request := new(models.GetShortRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	setDefaultRequest(request)
	sortAndOrder := request.SortBy + " " + request.Order
	offset := (request.Page - 1) * request.PageLength

	query := fmt.Sprintf(`
	SELECT shortlinks.Hash,shortlinks.OriginalUrl,domains.Name, shortlinks.Alias, shortlinks.CreatedOn, 
	shortlinks.ExpirationDate FROM shortlinks 
	JOIN domains on shortlinks.DomainId = domains.Id
	WHERE shortlinks.OrganizationId = %s AND shortlinks.IsDeprecated = %t AND domains.IsDeprecated = %t
	ORDER BY %s LIMIT %d OFFSET %d`,
		request.OrganizationId, false, false, sortAndOrder, request.PageLength, offset)

	rows, err := urlService.Db.Query(query)
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer rows.Close()

	var shorts []shortDto
	for rows.Next() {
		var short shortDto
		err = rows.Scan(&short.Hash, &short.OriginalUrl, &short.DomainName, &short.Alias, &short.CreatedOn, &short.ExpirationDate)
		if err != nil {
			return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		shorts = append(shorts, short)
	}

	query2 := `SELECT COUNT(1) FROM shortlinks WHERE OrganizationId =? AND IsDeprecated=?`
	row := urlService.Db.QueryRow(query2, request.OrganizationId, false)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	var shortDatas []models.GetShortData
	for _, short := range shorts {
		data := models.GetShortData{
			Short:          short.Hash,
			OriginalUrl:    short.OriginalUrl,
			Domain:         short.DomainName,
			Alias:          short.Alias,
			CreatedOn:      short.CreatedOn,
			ExpirationDate: short.ExpirationDate,
		}
		shortDatas = append(shortDatas, data)
	}

	totalPage := int(math.Ceil(float64(count) / float64(request.PageLength)))
	resp := models.GetShortResponse{
		ShortLinks: shortDatas,
		Page:       request.Page,
		TotalPage:  totalPage,
		Totals:     count,
		PageLength: len(shortDatas),
	}

	return urlContext.JSON(http.StatusOK, resp)
}

func setDefaultRequest(request *models.GetShortRequest) {
	if request.Page < 1 {
		request.Page = 1
	}

	if request.PageLength < 10 {
		request.PageLength = 10
	}

	if request.SortBy == "" {
		request.SortBy = "ExpirationDate"
	}

	if request.Order == "" {
		request.Order = "asc"
	}
}
