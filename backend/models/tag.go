package models

type CreateTagRequest struct {
	Name  string `json:"name"`
	Short string `json:"short"`
}

type CreateTagResponse struct {
	Id string `json:"id"`
}

type CreateShortLinkTagRequest struct {
	Short string `json:"short"`
	TagId string `json:"tagId"`
}

type CreateShortLinkTagResponse struct {
	Id string `json:"id"`
}
