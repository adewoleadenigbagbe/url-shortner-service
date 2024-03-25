package services

import (
	"database/sql"
	"errors"
	"net/http"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type UserService struct {
	Db *sql.DB
}

func (service UserService) SendEmail(userContext echo.Context) error {
	var err error
	request := new(models.SendEmailRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	errs := validateInvitesRequest(*request)
	if len(errs) != 0 {
		valErrors := lo.Map(errs, func(er error, index int) string {
			return er.Error()
		})
		return userContext.JSON(http.StatusBadRequest, valErrors)
	}

	recipients := lo.Map(request.Invites, func(invite models.Invite, index int) string {
		return invite.Email
	})

	err = helpers.SendEmail(recipients)
	if err != nil {
		return userContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	sqlStmt, args := formatInsertStatement(*request)
	_, err = service.Db.Exec(sqlStmt, args...)
	if err != nil {
		return userContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	return userContext.JSON(http.StatusOK, nil)
}

func validateInvitesRequest(request models.SendEmailRequest) []error {
	var validationErrors []error
	if request.ReferralId == "" {
		validationErrors = append(validationErrors, errors.New("referralId is required"))
	}

	if len(request.Invites) == 0 {
		validationErrors = append(validationErrors, errors.New("invites is empty. supply at least one"))
	}

	return validationErrors
}

func formatInsertStatement(request models.SendEmailRequest) (string, []interface{}) {
	stmt := "INSERT INTO invites VALUES "
	vals := []interface{}{}

	for _, invite := range request.Invites {
		stmt += "(?, ?, ?),"
		uuid := sequentialguid.NewSequentialGuid().String()
		vals = append(vals, uuid, invite.Email, request.ReferralId, invite.RoleId)
	}
	//trim the last ,
	stmt = stmt[0 : len(stmt)-1]
	return stmt, vals
}
