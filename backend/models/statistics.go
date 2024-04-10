package models

import (
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers/sqltype"
)

var _ IAggregateRow = (*CityAggregateRow)(nil)
var _ IAggregateRow = (*CountryAggregateRow)(nil)

type IAggregateRow interface {
	GetCount() int
}

type GetShortStatisticRequest struct {
	ShortId        string                      `query:"shortId"`
	OrganizationId string                      `header:"X-OrganizationId"`
	DateRangeType  enums.DateRange             `query:"dateRange"`
	StartDate      sqltype.Nullable[time.Time] `query:"startDate"`
	EndDate        sqltype.Nullable[time.Time] `query:"endDate"`
	Timezone       string                      `query:"timezone"`
	SortBy         string                      `query:"sortBy"`
}

type GetShortStatisticResponse struct {
	ShortId            string                 `json:"shortId"`
	Domain             string                 `json:"domain"`
	Hash               string                 `json:"hash"`
	CityAggregates     []CityAggregateRow     `json:"cityAggregates"`
	CountryAggregates  []CountryAggregateRow  `json:"countryAggregates"`
	OsAggregates       []OsAggregateRow       `json:"osAggregates"`
	PlatformAggregates []PlatformAggregateRow `json:"platformAggregates"`
	BrowserAggregates  []BrowserAggregateRow  `json:"browserAggregates"`
	DateAggregates     []DateAggregateRow     `json:"dateAggregates"`
}

type CityAggregateRow struct {
	City  string `json:"city"`
	Count int    `json:"count"`
}

func (cityAggregateRow CityAggregateRow) GetCount() int {
	return cityAggregateRow.Count
}

type CountryAggregateRow struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

type OsAggregateRow struct {
	Os    string `json:"os"`
	Count int    `json:"count"`
}

type PlatformAggregateRow struct {
	Platform string `json:"platform"`
	Count    int    `json:"count"`
}

type BrowserAggregateRow struct {
	Browser string `json:"browser"`
	Count   int    `json:"count"`
}

type DateAggregateRow struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func (countryAggregateRow CountryAggregateRow) GetCount() int {
	return countryAggregateRow.Count
}
