package models

import "database/sql"

type CreateUrlRequest struct {
	UserId      string         `json:"userId"`
	OriginalUrl string         `json:"original_url"`
	DomainName  string         `json:"domain"`
	CustomAlias sql.NullString `json:"alias"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"url"`
}
