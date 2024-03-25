package services

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (service UserService) ConvertReferral(userContext echo.Context) error {
	var err error
	request := new(models.ConvertReferralRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	//validate here
	errs := validateReferralRequest(*request)
	if len(errs) > 0 {
		valErrors := lo.Map(errs, func(er error, index int) string {
			return er.Error()
		})
		return userContext.JSON(http.StatusBadRequest, valErrors)
	}

	row := service.Db.QueryRow("SELECT OrganizationId FROM users WHERE Id=? AND IsDeprecated=?", request.ReferralId, false)
	var organizationId string
	err = row.Scan(&organizationId)
	if err != nil {
		return userContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	row2 := service.Db.QueryRow("SELECT RoleId FROM invites WHERE Email=?", request.Email)
	var roleId string
	err = row2.Scan(&roleId)
	if err != nil {
		return userContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if len(roleId) == 0 {
		return userContext.JSON(http.StatusBadRequest, []string{"invite does not exist"})
	}

	tx, err := service.Db.Begin()
	if err != nil {
		return userContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	userid := sequentialguid.NewSequentialGuid().String()
	hashedPassword := helpers.GeneratePassword(request.Password)
	usercreatedOn := time.Now()
	_, err = tx.Exec(`INSERT INTO users VALUES(?,?,?,?,?,?,?,?,?,?,?);`,
		userid, request.Username, request.Email, hashedPassword, usercreatedOn,
		usercreatedOn, usercreatedOn, roleId, organizationId, request.ReferralId, false)
	if err != nil {
		tx.Rollback()
		return userContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	//create userkeys
	userKeyId := sequentialguid.NewSequentialGuid().String()
	apikey := helpers.GenerateApiKey(request.Email)
	keyCreatedOn := time.Now()
	expiryDate := keyCreatedOn.AddDate(models.ApiExpiry, 0, 0)
	_, err = tx.Exec("INSERT INTO userkeys VALUES(?,?,?,?,?,?,?);", userKeyId, apikey, keyCreatedOn, keyCreatedOn, expiryDate, userid, true)
	if err != nil {
		tx.Rollback()
		return userContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	tx.Commit()
	return userContext.JSON(http.StatusCreated, models.ConvertReferralResponse{Id: userid, ApiKey: apikey})
}

func validateReferralRequest(request models.ConvertReferralRequest) []error {
	var validationErrors []error

	if request.Username == "" {
		validationErrors = append(validationErrors, errors.New("username is required"))
	}

	if request.ReferralId == "" {
		validationErrors = append(validationErrors, errors.New("referralId is required"))
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, errors.New("password is required"))
	}

	isEmailValid, _ := regexp.MatchString(models.EmailRegex, request.Email)
	if !isEmailValid {
		validationErrors = append(validationErrors, errors.New("email is invalid"))
	}
	return validationErrors
}
