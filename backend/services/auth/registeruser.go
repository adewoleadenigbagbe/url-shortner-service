package services

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type AuthService struct {
	Db  *sql.DB
	Rdb *redis.Client
}

func (service AuthService) RegisterUser(authContext echo.Context) error {
	var err error
	request := new(models.RegisterUserRequest)
	err = authContext.Bind(request)
	if err != nil {
		return authContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	//client validation
	errs := validateUser(*request)
	if len(errs) > 0 {
		valErrors := lo.Map(errs, func(er error, index int) string {
			return er.Error()
		})
		return authContext.JSON(http.StatusBadRequest, valErrors)
	}

	var roleId string
	row := service.Db.QueryRow("SELECT Id FROM userRoles WHERE Role=?", enums.Administrator)
	err = row.Scan(&roleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authContext.JSON(http.StatusNotFound, []string{"no role found"})
		}
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	var payPlanId string
	queryPayPlanRow := service.Db.QueryRow("SELECT Id FROM payplans WHERE Type=?", enums.Free)
	err = queryPayPlanRow.Scan(&payPlanId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authContext.JSON(http.StatusNotFound, []string{"payplan does not exist"})
		}
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	var count int
	queryOrgRow := service.Db.QueryRow("SELECT COUNT(1) FROM organizations WHERE Name=?", request.Company)
	err = queryOrgRow.Scan(&count)
	if err != nil {
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	if count > 0 {
		return authContext.JSON(http.StatusBadRequest, []string{"company name exist"})
	}

	tx, err := service.Db.Begin()
	if err != nil {
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	// create organization
	userid := sequentialguid.NewSequentialGuid().String()
	organizationId := sequentialguid.NewSequentialGuid().String()
	organizationCreatedOn := time.Now()
	_, err = tx.Exec(`INSERT INTO organizations VALUES(?,?,?,?,?,?,?,?);`,
		organizationId, request.Company, request.PhoneNumber, request.Timezone, userid, organizationCreatedOn,
		organizationCreatedOn, false)

	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	//organization payplan
	organizationPlanId := sequentialguid.NewSequentialGuid().String()
	planCreatedOn := time.Now()
	_, err = tx.Exec("INSERT INTO organizationpayplans VALUES(?,?,?,?,?,?,?);", organizationPlanId, enums.None, payPlanId, organizationId, planCreatedOn, planCreatedOn, enums.Current)
	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	//user
	hashedPassword := helpers.GeneratePassword(request.Password)
	usercreatedOn := time.Now()
	_, err = tx.Exec(`INSERT INTO users VALUES(?,?,?,?,?,?,?,?,?,?,?);`,
		userid, request.UserName, request.Email, hashedPassword, usercreatedOn,
		usercreatedOn, usercreatedOn, organizationId, sql.NullString{}, roleId, false)
	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	//userkeys
	userKeyId := sequentialguid.NewSequentialGuid().String()
	apikey := helpers.GenerateApiKey(request.Email)
	keyCreatedOn := time.Now()
	expiryDate := keyCreatedOn.AddDate(models.ApiExpiry, 0, 0)
	_, err = tx.Exec("INSERT INTO userkeys VALUES(?,?,?,?,?,?,?,?);", userKeyId, apikey, keyCreatedOn, keyCreatedOn, expiryDate, userid, organizationId, true)
	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	tx.Commit()
	return authContext.JSON(http.StatusOK, models.RegisterUserResponse{Id: userid, ApiKey: apikey})
}

func validateUser(user models.RegisterUserRequest) []error {
	var validationErrors []error

	if user.UserName == "" {
		validationErrors = append(validationErrors, errors.New("username is required"))
	}

	if user.Password == "" {
		validationErrors = append(validationErrors, errors.New("password is required"))
	}

	if user.Company == "" {
		validationErrors = append(validationErrors, errors.New("company name is required"))
	}

	if user.Timezone == "" {
		validationErrors = append(validationErrors, errors.New("timezone is required"))
	}

	isEmailValid, _ := regexp.MatchString(models.EmailRegex, user.Email)
	if !isEmailValid {
		validationErrors = append(validationErrors, errors.New("email is invalid"))
	}

	isPhoneValid, _ := regexp.MatchString(models.PhoneNumberRegex, user.PhoneNumber)
	if !isPhoneValid {
		validationErrors = append(validationErrors, errors.New("phone number is invalid"))
	}

	return validationErrors
}
