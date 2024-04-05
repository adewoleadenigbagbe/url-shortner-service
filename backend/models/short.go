package models

import (
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
)

type CreateUrlRequest struct {
	OriginalUrl    string                   `json:"originalurl"`
	DomainId       string                   `json:"domainId"`
	CustomAlias    helpers.Nullable[string] `json:"alias"`
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
	Country        helpers.Nullable[string] `json:"country"`
	TimeZone       helpers.Nullable[string] `json:"timezone"`
	City           helpers.Nullable[string] `json:"city"`
	Os             helpers.Nullable[string] `json:"os"`
	Browser        helpers.Nullable[string] `json:"browser"`
	UserAgent      helpers.Nullable[string] `json:"userAgent"`
	Platform       helpers.Nullable[string] `json:"platform"`
	IpAddress      helpers.Nullable[string] `json:"ipAddress"`
	Method         helpers.Nullable[string] `json:"method"`
	Status         helpers.Nullable[int]    `json:"status"`
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
	Alias          helpers.Nullable[string] `json:"alias"`
	CreatedOn      time.Time                `json:"createdOn"`
	ExpirationDate time.Time                `json:"expiryDate"`
}
