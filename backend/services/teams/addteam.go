package services

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

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
	binder := &echo.DefaultBinder{}

	err = binder.BindHeaders(teamContext, request)
	if err != nil {
		return teamContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	err = binder.BindBody(teamContext, request)
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

	queryStr, args := formatQuery(*request)
	rows, err := service.Db.Query(queryStr, args...)
	if err != nil {
		return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	defer rows.Close()

	var teams []string
	for rows.Next() {
		var team string
		err = rows.Scan(&team)
		if err != nil {
			return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		teams = append(teams, team)
	}

	if len(teams) > 0 {
		errs := lo.Map(teams, func(team string, index int) string {
			return team + " already exist"
		})
		return teamContext.JSON(http.StatusBadRequest, errs)
	}

	insertStmt, insertArgs := formatTeamStmt(*request)
	_, err = service.Db.Exec(insertStmt, insertArgs...)
	if err != nil {
		return teamContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	return teamContext.JSON(http.StatusCreated, nil)
}

func validateAddTeamRequest(request models.CreateTeamRequest) []error {
	var validationErrors []error
	if len(request.Teams) == 0 {
		validationErrors = append(validationErrors, errors.New("teams is empty. supply at least one"))
	}

	if request.OrganizationId == "" || len(request.OrganizationId) < 36 {
		validationErrors = append(validationErrors, errors.New("organizationId is required"))
	}

	return validationErrors
}

// SELECT Name FROM teams WHERE Name IN (?,?,?) AND OrganizationId =? AND IsDeprecated = false
func formatQuery(request models.CreateTeamRequest) (string, []interface{}) {
	str := "SELECT Name FROM teams WHERE Name IN (?" + strings.Repeat(",?", len(request.Teams)-1) + ") AND OrganizationId =? AND IsDeprecated =?"
	vals := []interface{}{}
	for _, team := range request.Teams {
		vals = append(vals, team)
	}
	vals = append(vals, request.OrganizationId, false)
	return str, vals
}

func formatTeamStmt(request models.CreateTeamRequest) (string, []interface{}) {
	stmt := "INSERT INTO teams VALUES "
	vals := []interface{}{}

	for _, team := range request.Teams {
		stmt += "(?,?,?,?),"
		uuid := sequentialguid.NewSequentialGuid().String()
		vals = append(vals, uuid, team, request.OrganizationId, false)
	}
	//trim the last ,
	stmt = stmt[0 : len(stmt)-1]
	return stmt, vals
}
