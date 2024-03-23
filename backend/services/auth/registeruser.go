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

const (
	expiryYear       = 1
	emailRegex       = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
	PhoneNumberRegex = "\\+[1-9]{1}[0-9]{0,2}-[2-9]{1}[0-9]{2}-[2-9]{1}[0-9]{2}-[0-9]{4}$"
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

	row := service.Db.QueryRow("SELECT Id FROM userRoles WHERE Role=?", enums.Administrator)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return authContext.JSON(http.StatusNotFound, []string{"no role found"})
	}

	var roleId string
	err = row.Scan(&roleId)
	if err != nil {
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
	_, err = tx.Exec(`INSERT INTO organizations VALUES(?,?,?,?,?,?,?);`,
		organizationId, request.Company, request.PhoneNumber, userid, organizationCreatedOn,
		organizationCreatedOn, false)

	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	//create user
	hashedPassword := helpers.GeneratePassword(request.Password)
	usercreatedOn := time.Now()
	_, err = tx.Exec(`INSERT INTO users VALUES(?,?,?,?,?,?,?,?,?,?,?);`,
		userid, request.UserName, request.Email, hashedPassword, usercreatedOn,
		usercreatedOn, usercreatedOn, roleId, organizationId, sql.NullString{}, false)
	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	//create userkeys
	userKeyId := sequentialguid.NewSequentialGuid().String()
	apikey := helpers.GenerateApiKey(request.Email)
	keyCreatedOn := time.Now()
	expiryDate := keyCreatedOn.AddDate(expiryYear, 0, 0)
	_, err = tx.Exec("INSERT INTO userkeys VALUES(?,?,?,?,?,?,?);", userKeyId, apikey, keyCreatedOn, keyCreatedOn, expiryDate, userid, true)
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

	isEmailValid, _ := regexp.MatchString(emailRegex, user.Email)
	if !isEmailValid {
		validationErrors = append(validationErrors, errors.New("email is invalid"))
	}

	isPhoneValid, _ := regexp.MatchString(PhoneNumberRegex, user.PhoneNumber)
	if !isPhoneValid {
		validationErrors = append(validationErrors, errors.New("phone number is invalid"))
	}

	return validationErrors
}
