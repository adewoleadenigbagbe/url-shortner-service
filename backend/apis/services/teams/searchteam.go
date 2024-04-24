package services

import (
	"fmt"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service TeamService) SearchTeam(teamContext echo.Context) error {
	var err error
	request := new(models.SearchTeamRequest)
	binder := &echo.DefaultBinder{}

	err = binder.BindHeaders(teamContext, request)
	if err != nil {
		return teamContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	err = binder.BindQueryParams(teamContext, request)
	if err != nil {
		return teamContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	term := "'%" + request.Term + "%'"
	query := fmt.Sprintf("SELECT Id,Name FROM teams WHERE OrganizationId = '%s' AND IsDeprecated = %t AND Name LIKE %s", request.OrganizationId, false, term)
	rows, err := service.Db.Query(query)
	if err != nil {
		return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	defer rows.Close()

	var teams []models.TeamDTO
	for rows.Next() {
		var team models.TeamDTO
		err = rows.Scan(&team.Id, &team.Name)
		if err != nil {
			return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		teams = append(teams, team)
	}

	return teamContext.JSON(http.StatusOK, models.SearchTeamResponse{Teams: teams})
}
