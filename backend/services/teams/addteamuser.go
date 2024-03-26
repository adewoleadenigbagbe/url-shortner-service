package services

import (
	"errors"
	"net/http"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (service TeamService) AddUserToTeam(teamContext echo.Context) error {
	var err error
	request := new(models.CreateUserTeamRequest)
	err = teamContext.Bind(request)
	if err != nil {
		return teamContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	modelErrs := validateAddUserTeamRequest(*request)
	if len(modelErrs) > 0 {
		valErrors := lo.Map(modelErrs, func(er error, index int) string {
			return er.Error()
		})
		return teamContext.JSON(http.StatusBadRequest, valErrors)
	}

	id := sequentialguid.NewSequentialGuid().String()
	_, err = service.Db.Exec("INSERT INTO teamusers VALUES(?,?,?);", id, request.TeamId, request.UserId)
	if err != nil {
		return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	return teamContext.JSON(http.StatusCreated, models.CreateUserTeamResponse{Id: id})
}

func validateAddUserTeamRequest(request models.CreateUserTeamRequest) []error {
	var validationErrors []error
	if request.TeamId == "" || len(request.TeamId) < 36 {
		validationErrors = append(validationErrors, errors.New("teamId is required"))
	}

	if request.UserId == "" || len(request.UserId) < 36 {
		validationErrors = append(validationErrors, errors.New("userId is required"))
	}

	return validationErrors
}
