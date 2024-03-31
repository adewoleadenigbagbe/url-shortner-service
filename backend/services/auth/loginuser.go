package services

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service AuthService) LoginUser(authContext echo.Context) error {
	var err error
	request := new(models.SignInUserRequest)
	err = authContext.Bind(request)

	if err != nil {
		return authContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	var apikey string
	var id string
	var role enums.Role
	query := `SELECT users.Id, userkeys.ApiKey, userRoles.Role FROM users 
					JOIN userkeys ON users.Id = userkeys.UserId
					JOIN userRoles ON users.RoleId = userRoles.Id
					WHERE users.Email=?
					AND users.IsDeprecated=?
					AND userRoles.IsDeprecated=? 
					AND userkeys.IsActive=?`

	row := service.Db.QueryRow(query, request.Email, false, false, true)
	if err = row.Scan(&id, &apikey, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authContext.JSON(http.StatusBadRequest, []string{"email is incorrect"})
		} else {
			return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
	}

	request.Id = id
	request.Role = role
	token, err := helpers.GenerateJWT(request.Id, request.Role, request.Email)
	if err != nil {
		return authContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	now := time.Now()
	service.Db.Exec("UPDATE users SET LastLogin =? , ModifiedOn=? WHERE Id=?", now, now, id)
	return authContext.JSON(http.StatusOK, models.SignInUserResponse{Token: token, Id: id, ApiKey: apikey})
}
