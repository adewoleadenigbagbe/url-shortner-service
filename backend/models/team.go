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
