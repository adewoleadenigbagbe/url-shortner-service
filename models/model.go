package models

type CreateUrlRequest struct {
	UserId      string `json:"userId"`
	OriginalUrl string `json:"original_url"`
	DomainName  string `json:"domain"`
	CustomAlias string `json:"alias"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"url"`
}
