package models

import (
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers/sqltype"
)

type CreateUrlRequest struct {
	OriginalUrl    string                   `json:"originalurl"`
	DomainId       string                   `json:"domainId"`
	CustomAlias    sqltype.Nullable[string] `json:"alias"`
	UserId         string                   `json:"userId"`
	Cloaking       bool                     `json:"cloaking"`
	OrganizationId string                   `header:"X-OrganizationId"`
}

type CreateUrlResponse struct {
	Id       string `json:"id"`
	ShortUrl string `json:"shortlink"`
	DomainId string `json:"domainId"`
}

type DeleteUrlRequest struct {
	ShortUrl string `json:"url"`
}

type RedirectShortRequest struct {
	ShortUrl       string                   `json:"shortUrl"`
	Country        sqltype.Nullable[string] `json:"country"`
	TimeZone       sqltype.Nullable[string] `json:"timezone"`
	City           sqltype.Nullable[string] `json:"city"`
	Os             sqltype.Nullable[string] `json:"os"`
	Browser        sqltype.Nullable[string] `json:"browser"`
	UserAgent      sqltype.Nullable[string] `json:"userAgent"`
	Platform       sqltype.Nullable[string] `json:"platform"`
	IpAddress      sqltype.Nullable[string] `json:"ipAddress"`
	Method         sqltype.Nullable[string] `json:"method"`
	Status         sqltype.Nullable[int]    `json:"status"`
	OrganizationId string                   `json:"organizationId"`
}

type GetShortRequest struct {
	OrganizationId string `header:"X-OrganizationId"`
	Page           int    `query:"page"`
	PageLength     int    `query:"pageLength"`
	SortBy         string `query:"sortBy"`
	Order          string `query:"order"`
}

type GetShortResponse struct {
	ShortLinks []GetShortData `json:"shorts"`
	Page       int            `json:"page"`
	TotalPage  int            `json:"totalPage"`
	Totals     int            `json:"totals"`
	PageLength int            `json:"pageLength"`
}

type GetShortData struct {
	Short          string                   `json:"short"`
	OriginalUrl    string                   `json:"originalUrl"`
	Domain         string                   `json:"domain"`
	Alias          sqltype.Nullable[string] `json:"alias"`
	CreatedOn      time.Time                `json:"createdOn"`
	ExpirationDate time.Time                `json:"expiryDate"`
}
