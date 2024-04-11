package statistics

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

type StatisticsService struct {
	Db *sql.DB
}

func (service StatisticsService) GetShortStatistics(statisticsContext echo.Context) error {
	var err error
	request := new(models.GetShortStatisticRequest)
	binder := &echo.DefaultBinder{}
	err = binder.BindHeaders(statisticsContext, request)
	if err != nil {
		return statisticsContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	err = binder.BindQueryParams(statisticsContext, request)
	if err != nil {
		return statisticsContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	shortlinkRow := service.Db.QueryRow(`SELECT shortlinks.Hash, domains.Name FROM shortlinks
	JOIN domains ON shortlinks.DomainId = domains.Id
	WHERE shortlinks.Id =?
	AND shortlinks.OrganizationId =?
	AND shortlinks.IsDeprecated =?
	AND domains.IsDeprecated =?`,
		request.ShortId, request.OrganizationId, false, false)

	var link string
	var domain string
	err = shortlinkRow.Scan(&link, &domain)
	if err != nil {
		return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	startDate, endDate := request.DateRangeType.GetRanges(request.StartDate, request.EndDate)
	if startDate.IsZero() || endDate.IsZero() {
		return statisticsContext.JSON(http.StatusBadRequest, []string{"select a range date type"})
	}

	startDate = helpers.StartOfDay(startDate)
	endDate = helpers.EndOfDay(endDate)

	cityQuery := fmt.Sprintf(`
	SELECT accesslogs.City, COUNT(accesslogs.City) AS CityCount FROM accesslogs
	JOIN shortlinks ON accesslogs.ShortId = shortlinks.Id
	WHERE accesslogs.ShortId = '%s'
	AND accesslogs.OrganizationId = '%s'
	AND accesslogs.TimeZone = '%s'
	AND accesslogs.CreatedOn >= '%s' AND accesslogs.CreatedOn <= '%s'
	AND accesslogs.IsDeprecated = %t
	AND shortlinks.IsDeprecated = %t
	GROUP BY accesslogs.City
	ORDER BY CityCount %s
	`, request.ShortId, request.OrganizationId, request.Timezone, startDate, endDate, false, false, request.SortBy)

	cityrows, err := service.Db.Query(cityQuery)
	if err != nil {
		return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer cityrows.Close()

	var cityAggregates []models.CityAggregateRow
	for cityrows.Next() {
		var cityAggregate models.CityAggregateRow
		err = cityrows.Scan(&cityAggregate.City, &cityAggregate.Count)
		if err != nil {
			return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		cityAggregates = append(cityAggregates, cityAggregate)
	}

	countryQuery := fmt.Sprintf(`
	SELECT accesslogs.Country, COUNT(accesslogs.Country) AS CountryCount FROM accesslogs
	JOIN shortlinks ON accesslogs.ShortId = shortlinks.Id
	WHERE accesslogs.ShortId = '%s'
	AND accesslogs.OrganizationId = '%s'
	AND accesslogs.TimeZone = '%s'
	AND accesslogs.CreatedOn >= '%s' AND accesslogs.CreatedOn <= '%s'
	AND accesslogs.IsDeprecated = %t
	AND shortlinks.IsDeprecated = %t
	GROUP BY accesslogs.Country
	ORDER BY CountryCount %s
	`, request.ShortId, request.OrganizationId, request.Timezone, startDate, endDate, false, false, request.SortBy)

	countryrows, err := service.Db.Query(countryQuery)
	if err != nil {
		return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer countryrows.Close()

	var countryAggregates []models.CountryAggregateRow

	for countryrows.Next() {
		var countryAggregate models.CountryAggregateRow
		err = countryrows.Scan(&countryAggregate.Country, &countryAggregate.Count)
		if err != nil {
			return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		countryAggregates = append(countryAggregates, countryAggregate)
	}

	osQuery := fmt.Sprintf(`
	SELECT accesslogs.Os, COUNT(accesslogs.Os) AS OsCount FROM accesslogs
	JOIN shortlinks ON accesslogs.ShortId = shortlinks.Id
	WHERE accesslogs.ShortId = '%s'
	AND accesslogs.OrganizationId = '%s'
	AND accesslogs.TimeZone = '%s'
	AND accesslogs.CreatedOn >= '%s' AND accesslogs.CreatedOn <= '%s'
	AND accesslogs.IsDeprecated = %t
	AND shortlinks.IsDeprecated = %t
	GROUP BY accesslogs.Os
	ORDER BY OsCount %s
	`, request.ShortId, request.OrganizationId, request.Timezone, startDate, endDate, false, false, request.SortBy)

	osrows, err := service.Db.Query(osQuery)
	if err != nil {
		return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer osrows.Close()

	var osAggregates []models.OsAggregateRow

	for osrows.Next() {
		var osAggregate models.OsAggregateRow
		err = osrows.Scan(&osAggregate.Os, &osAggregate.Count)
		if err != nil {
			return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		osAggregates = append(osAggregates, osAggregate)
	}

	platformQuery := fmt.Sprintf(`
	SELECT accesslogs.Platform, COUNT(accesslogs.Platform) AS PlatformCount FROM accesslogs
	JOIN shortlinks ON accesslogs.ShortId = shortlinks.Id
	WHERE accesslogs.ShortId = '%s'
	AND accesslogs.OrganizationId = '%s'
	AND accesslogs.TimeZone = '%s'
	AND accesslogs.CreatedOn >= '%s' AND accesslogs.CreatedOn <= '%s'
	AND accesslogs.IsDeprecated = %t
	AND shortlinks.IsDeprecated = %t
	GROUP BY accesslogs.Platform
	ORDER BY PlatformCount %s
	`, request.ShortId, request.OrganizationId, request.Timezone, startDate, endDate, false, false, request.SortBy)

	platformrows, err := service.Db.Query(platformQuery)
	if err != nil {
		return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer platformrows.Close()

	var platformAggregates []models.PlatformAggregateRow

	for platformrows.Next() {
		var platformAggregate models.PlatformAggregateRow
		err = platformrows.Scan(&platformAggregate.Platform, &platformAggregate.Count)
		if err != nil {
			return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		platformAggregates = append(platformAggregates, platformAggregate)
	}

	browserQuery := fmt.Sprintf(`
	SELECT accesslogs.Browser, COUNT(accesslogs.Browser) AS BrowserCount FROM accesslogs
	JOIN shortlinks ON accesslogs.ShortId = shortlinks.Id
	WHERE accesslogs.ShortId = '%s'
	AND accesslogs.OrganizationId = '%s'
	AND accesslogs.TimeZone = '%s'
	AND accesslogs.CreatedOn >= '%s' AND accesslogs.CreatedOn <= '%s'
	AND accesslogs.IsDeprecated = %t
	AND shortlinks.IsDeprecated = %t
	GROUP BY accesslogs.Browser
	ORDER BY BrowserCount %s
	`, request.ShortId, request.OrganizationId, request.Timezone, startDate, endDate, false, false, request.SortBy)

	browserrows, err := service.Db.Query(browserQuery)
	if err != nil {
		return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer browserrows.Close()

	var browserAggregates []models.BrowserAggregateRow

	for browserrows.Next() {
		var browserAggregate models.BrowserAggregateRow
		err = browserrows.Scan(&browserAggregate.Browser, &browserAggregate.Count)
		if err != nil {
			return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		browserAggregates = append(browserAggregates, browserAggregate)
	}

	dateQuery := fmt.Sprintf(`
	SELECT Date(accesslogs.CreatedOn) AS Date, COUNT(accesslogs.CreatedOn) AS CreatedOnCount FROM accesslogs
	JOIN shortlinks ON accesslogs.ShortId = shortlinks.Id
	WHERE accesslogs.ShortId = '%s'
	AND accesslogs.OrganizationId = '%s'
	AND accesslogs.TimeZone = '%s'
	AND accesslogs.CreatedOn >= '%s' AND accesslogs.CreatedOn <= '%s'
	AND accesslogs.IsDeprecated = %t
	AND shortlinks.IsDeprecated = %t
	GROUP BY Date
	ORDER BY CreatedOnCount %s
	`, request.ShortId, request.OrganizationId, request.Timezone, startDate, endDate, false, false, request.SortBy)

	daterrows, err := service.Db.Query(dateQuery)
	if err != nil {
		return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer daterrows.Close()

	var dateAggregates []models.DateAggregateRow

	for daterrows.Next() {
		var dateAggregate models.DateAggregateRow
		err = daterrows.Scan(&dateAggregate.Date, &dateAggregate.Count)
		if err != nil {
			return statisticsContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		dateAggregates = append(dateAggregates, dateAggregate)
	}

	resp := models.GetShortStatisticResponse{
		ShortId:            request.ShortId,
		Domain:             domain,
		Hash:               link,
		CityAggregates:     cityAggregates,
		CountryAggregates:  countryAggregates,
		OsAggregates:       osAggregates,
		PlatformAggregates: platformAggregates,
		BrowserAggregates:  browserAggregates,
		DateAggregates:     dateAggregates,
	}
	return statisticsContext.JSON(http.StatusOK, resp)
}

// func GetAggregateRow(schema, column string, startDate, endDate time.Time, db *sql.DB, request models.GetShortStatisticRequest) ([]models.IAggregateRow, error) {
// 	columnField := schema + "." + column
// 	query := fmt.Sprintf(`
// 	SELECT %s , COUNT(%s) AS CountryCount FROM accesslogs
// 	JOIN shortlinks ON accesslogs.ShortId = shortlinks.Id
// 	WHERE accesslogs.ShortId = '%s'
// 	AND accesslogs.OrganizationId = '%s'
// 	AND accesslogs.TimeZone = '%s'
// 	AND accesslogs.CreatedOn >= '%s' AND accesslogs.CreatedOn <= '%s'
// 	AND accesslogs.IsDeprecated = %t
// 	AND shortlinks.IsDeprecated = %t
// 	GROUP BY %s
// 	ORDER BY CountryCount %s
// 	`, columnField, columnField, request.ShortId, request.OrganizationId, request.Timezone, startDate, endDate, false, false, columnField, request.SortBy)

// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var rowAggregates []models.IAggregateRow
// 	for rows.Next() {
// 		var rowAggregate models.IAggregateRow
// 		err = rows.Scan(rowAggregate)
// 		if err != nil {
// 			return nil, err
// 		}
// 		rowAggregates = append(rowAggregates, rowAggregate)
// 	}

// 	return rowAggregates, nil
// }
