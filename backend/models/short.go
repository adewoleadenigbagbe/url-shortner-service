package models

import (
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
)

type CreateUrlRequest struct {
	UserId      string                   `json:"userId"`
	OriginalUrl string                   `json:"originalurl"`
	DomainName  string                   `json:"domain"`
	CustomAlias helpers.Nullable[string] `json:"alias"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"shortlink"`
	Domain   string `json:"domain"`
}

type DeleteUrlRequest struct {
	ShortUrl string `json:"url"`
}

type RedirectShortRequest struct {
	ShortUrl string `json:"shorturl"`
}

type GetShortRequest struct {
	UserId     string `json:"userId"`
	Page       int    `query:"page"`
	PageLength int    `query:"pageLength"`
	SortBy     string `query:"sortBy"`
	Order      string `query:"order"`
}

type GetShortResponse struct {
	ShortLinks []GetShortData `json:"shorts"`
	Page       int            `json:"page"`
	TotalPage  int            `json:"totalPage"`
	Totals     int            `json:"totals"`
	PageLength int            `json:"pageLength"`
}

type GetShortData struct {
	Short          string    `json:"short"`
	OriginalUrl    string    `json:"originalUrl"`
	Domain         string    `json:"domain"`
	Alias          string    `json:"alias"`
	CreatedOn      time.Time `json:"createdOn"`
	ExpirationDate time.Time `json:"expiryDate"`
	UserId         string    `json:"userId"`
}
