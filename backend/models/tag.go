package models

type CreateTagRequest struct {
	Tags []string `json:"tags"`
}

type CreateTagResponse struct {
	Id string `json:"id"`
}

type CreateShortLinkTagRequest struct {
	Tags    []string `json:"tags"`
	ShortId string   `json:"shortId"`
}

type CreateShortLinkTagResponse struct {
	Id string `json:"id"`
}

type SearchTagRequest struct {
	Term string `query:"term"`
}

type SearchTagResponse struct {
	Tags []string `json:"tags"`
}
