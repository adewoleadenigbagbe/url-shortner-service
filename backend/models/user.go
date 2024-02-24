package models

type RegisterUserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserResponse struct {
	Id     string `json:"userId"`
	ApiKey string `json:"apikey"`
}

type SigInUserRequest struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type SignInUserResponse struct {
	Id     string `json:"userId"`
	ApiKey string `json:"apikey"`
	Token  string `json:"access_token"`
}

type GetShortRequest struct {
	ShortUrl string `json:"shorturl"`
}
