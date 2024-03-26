package services

import (
	"database/sql"
	"errors"
	"net/http"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type TeamService struct {
	Db *sql.DB
}

func (service TeamService) AddTeam(teamContext echo.Context) error {
	var err error
	request := new(models.CreateTeamRequest)
	err = teamContext.Bind(request)
	if err != nil {
		return teamContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	modelErrs := validateAddTeamRequest(*request)
	if len(modelErrs) > 0 {
		valErrors := lo.Map(modelErrs, func(er error, index int) string {
			return er.Error()
		})
		return teamContext.JSON(http.StatusBadRequest, valErrors)
	}

	var orgName string
	row := service.Db.QueryRow("SELECT Name FROM teams WHERE organizationId =? AND IsDeprecated=?", request.OrganizationId, false)
	err = row.Scan(&orgName)
	if err != nil {
		return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if len(orgName) > 0 {
		return teamContext.JSON(http.StatusBadRequest, []string{"team name already exist"})
	}

	id := sequentialguid.NewSequentialGuid().String()
	_, err = service.Db.Exec("INSERT INTO teams VALUES(?,?,?,?);", id, request.Name, request.OrganizationId, false)
	if err != nil {
		return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	return teamContext.JSON(http.StatusCreated, models.CreateTeamResponse{Id: id})
}

func validateAddTeamRequest(request models.CreateTeamRequest) []error {
	var validationErrors []error
	if request.Name == "" {
		validationErrors = append(validationErrors, errors.New("name is required"))
	}

	if request.OrganizationId == "" || len(request.OrganizationId) < 36 {
		validationErrors = append(validationErrors, errors.New("organizationId is required"))
	}

	return validationErrors
}
