package models

import "github.com/adewoleadenigbagbe/url-shortner-service/enums"

type RegisterUserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserResponse struct {
	Id     string `json:"userId"`
	ApiKey string `json:"apikey"`
}

type SignInUserRequest struct {
	Email string `json:"email"`
	Id    string
	Role  enums.Role
}

type SignInUserResponse struct {
	Id     string `json:"userId"`
	ApiKey string `json:"apikey"`
	Token  string `json:"access_token"`
}

type SignOutUserRequest struct {
	UserId string `json:"userId"`
}
