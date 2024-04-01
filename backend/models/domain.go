package models

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
