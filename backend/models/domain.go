package models

type CreateDomainRequest struct {
	Name   string `json:"domain"`
	UserId string `json:"userId"`
}

type DeleteDomainRequest struct {
	Name string `json:"domain"`
}
