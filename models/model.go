package models

type CreateUrlRequest struct {
	UserId      string `json:"userId"`
	OriginalUrl string `json:"original_url"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"url"`
}
