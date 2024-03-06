package models

import (
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
