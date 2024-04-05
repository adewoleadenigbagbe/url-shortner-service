package models

import "time"

type CreateDomainRequest struct {
	Name           string `json:"domain"`
	UserId         string `json:"userId"`
	IsCustom       bool   `json:"isCustom"`
	OrganizationId string `header:"X-OrganizationId"`
}

type CreateDomainResponse struct {
	DomainId string `json:"domainId"`
	Name     string `json:"name"`
}

type DeleteDomainRequest struct {
	Name string `json:"domain"`
}

type GetDomainRequest struct {
	OrganizationId string `header:"userId"`
	Page           int    `query:"page"`
	PageLength     int    `query:"pageLength"`
	SortBy         string `query:"sortBy"`
	Order          string `query:"order"`
}

type GetDomainData struct {
	Name      string    `json:"name"`
	IsCustom  bool      `json:"isCustom"`
	CreatedOn time.Time `json:"createdOn"`
}

type GetDomainResponse struct {
	Domains    []GetDomainData `json:"domains"`
	Page       int             `json:"page"`
	TotalPage  int             `json:"totalPage"`
	Totals     int             `json:"totals"`
	PageLength int             `json:"pageLength"`
}
