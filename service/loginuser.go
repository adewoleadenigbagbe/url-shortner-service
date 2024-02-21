package services

import (
	"database/sql"
	"net/http"

	jwtauth "github.com/adewoleadenigbagbe/url-shortner-service/helpers/auth"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service AuthService) LoginUser(userContext echo.Context) error {
	var err error
	request := new(models.SigInUserRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	var apikey string
	var id string

	query := "SELECT users.Id , userkeys.ApiKey FROM users JOIN userkeys ON users.Id = userkeys.UserId WHERE users.Email=? AND userkeys.IsActive=?"
	row := service.Db.QueryRow(query, request.Email, true)
	if err = row.Scan(&id, &apikey); err == sql.ErrNoRows {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	token, err := jwtauth.GenerateJWT(*request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	return userContext.JSON(http.StatusOK, models.SignInUserResponse{Token: token, Id: id, ApiKey: apikey})
}
