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
	expiryYear = 1
	emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
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
		return authContext.JSON(http.StatusBadRequest, err.Error())
	}

	//client validation
	errs := validateUser(*request)
	if len(errs) > 0 {
		valErrors := lo.Map(errs, func(er error, index int) string {
			return er.Error()
		})
		return authContext.JSON(http.StatusBadRequest, valErrors)
	}

	//generate api key
	userid := sequentialguid.NewSequentialGuid().String()
	usercreatedOn := time.Now()

	row := service.Db.QueryRow("SELECT Id FROM userRoles WHERE Role=?", enums.User)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return authContext.JSON(http.StatusNotFound, errors.New("no role found"))
	}
	var roleId string
	err = row.Scan(&roleId)
	if err != nil {
		return authContext.JSON(http.StatusInternalServerError, err.Error())
	}

	tx, err := service.Db.Begin()
	if err != nil {
		return authContext.JSON(http.StatusInternalServerError, err.Error())
	}

	//save user
	_, err = tx.Exec(`INSERT INTO users VALUES(?,?,?,?,?,?,?,?);`,
		userid, request.UserName, request.Email, usercreatedOn,
		usercreatedOn, usercreatedOn, roleId, false)

	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, err.Error())
	}

	//save keys
	userKeyId := sequentialguid.NewSequentialGuid().String()
	apikey := helpers.GenerateApiKey(request.Email)
	keyCreatedOn := time.Now()
	expiryDate := keyCreatedOn.AddDate(expiryYear, 0, 0)

	_, err = tx.Exec("INSERT INTO userkeys VALUES(?,?,?,?,?,?,?);", userKeyId, apikey, keyCreatedOn, keyCreatedOn, expiryDate, userid, true)
	if err != nil {
		tx.Rollback()
		return authContext.JSON(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()
	return authContext.JSON(http.StatusOK, models.RegisterUserResponse{Id: userid, ApiKey: apikey})
}

func validateUser(user models.RegisterUserRequest) []error {
	var validationErrors []error

	if user.UserName == "" {
		validationErrors = append(validationErrors, errors.New("username is required"))
	}

	isEmailValid, _ := regexp.MatchString(emailRegex, user.Email)
	if !isEmailValid {
		validationErrors = append(validationErrors, errors.New("email is invalid"))
	}
	return validationErrors
}
