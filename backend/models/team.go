package models

type CreateTeamRequest struct {
	Name           string `json:"name"`
	OrganizationId string `json:"organizationId"`
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
