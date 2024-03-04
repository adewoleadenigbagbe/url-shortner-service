package services

import (
	"database/sql"
	"errors"
	"net/http"

	jwtauth "github.com/adewoleadenigbagbe/url-shortner-service/helpers/auth"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service AuthService) LoginUser(authContext echo.Context) error {
	var err error
	request := new(models.SignInUserRequest)
	err = authContext.Bind(request)

	if err != nil {
		return authContext.JSON(http.StatusBadRequest, err.Error())
	}

	var apikey string
	var id string
	query := "SELECT users.Id, userkeys.ApiKey FROM users JOIN userkeys ON users.Id = userkeys.UserId WHERE users.Email=? AND userkeys.IsActive=?"
	row := service.Db.QueryRow(query, request.Email, true)
	if err = row.Scan(&id, &apikey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authContext.JSON(http.StatusBadRequest, errors.New("email is incorrect").Error())
		} else {
			return authContext.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	request.Id = id
	token, err := jwtauth.GenerateJWT(*request)
	if err != nil {
		return authContext.JSON(http.StatusBadRequest, err.Error())
	}

	return authContext.JSON(http.StatusOK, models.SignInUserResponse{Token: token, Id: id, ApiKey: apikey})
}
