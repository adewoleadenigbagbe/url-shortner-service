package models

type UserRequest struct {
	UserName string
	Email    string
}

type UserResponse struct {
	Id     int
	ApiKey string
}
