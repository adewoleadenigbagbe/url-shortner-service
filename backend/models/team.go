package models

type CreateTeamRequest struct {
	Teams          []string `json:"teams"`
	OrganizationId string   `header:"X-OrganizationId"`
}

type CreateTeamResponse struct {
	Id string `json:"id"`
}

type CreateUserTeamRequest struct {
	TeamId string `json:"teamId"`
	UserId string `json:"userId"`
}

type CreateUserTeamResponse struct {
	Id string `json:"id"`
}

type SearchTeamRequest struct {
	Term           string `query:"term"`
	OrganizationId string `header:"X-OrganizationId"`
}

type SearchTeamResponse struct {
	Teams []TeamDTO `json:"teams"`
}

type TeamDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
