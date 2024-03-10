package services

import (
	"fmt"
	"math"
	"net/http"
	"time"

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
		DomainId       string
		Alias          string
		CreatedOn      time.Time
		ExpirationDate time.Time
	}

	request := new(models.GetShortRequest)
	err = urlContext.Bind(request)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, err.Error())
	}

	setDefaultRequest(request)

	sortAndOrder := request.SortBy + " " + request.Order
	query := `SELECT Hash,OriginalUrl,DomainId, Alias, CreatedOn, ExpirationDate 
	FROM shortlinks 
	WHERE userId =? AND IsDepecated=?
	ORDER BY ?=
	LIMIT =? OFFSET =?`

	rows, err := urlService.Db.Query(query, request.UserId, false, sortAndOrder, request.PageLength, request.Page)
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var shorts []shortDto
	for rows.Next() {
		var short shortDto
		err = rows.Scan(&short)
		if err != nil {
			return urlContext.JSON(http.StatusInternalServerError, err.Error())
		}
		shorts = append(shorts, short)
	}
	fmt.Println(shorts)

	query2 := `SELECT COUNT(1) FROM shortlinks WHERE userId =? AND IsDepecated=?`
	row := urlService.Db.QueryRow(query2, request.UserId, false)
	if row.Err() != nil {
		return urlContext.JSON(http.StatusInternalServerError, row.Err().Error())
	}

	var count int
	err = row.Scan(&count)
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, err)
	}

	var shortDatas []models.GetShortData
	for _, short := range shorts {
		data := models.GetShortData{
			Short:          short.Hash,
			OriginalUrl:    short.OriginalUrl,
			Domain:         short.DomainId,
			Alias:          short.Alias,
			CreatedOn:      short.CreatedOn,
			ExpirationDate: short.ExpirationDate,
			UserId:         request.UserId,
		}
		shortDatas = append(shortDatas, data)
	}

	totalPage := int(math.Ceil(float64(count) / float64(request.PageLength)))
	resp := models.GetShortResponse{
		ShortLinks: shortDatas,
		Page:       request.Page,
		TotalPage:  totalPage,
		Totals:     count,
		PageLength: request.PageLength,
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
		request.SortBy = "expirationDate"
	}

	if request.Order == "" {
		request.Order = "asc"
	}
}
