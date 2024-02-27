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
)

const (
	expiryYear = 1
	emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
)

func (service AuthService) RegisterUser(userContext echo.Context) error {
	var err error
	request := new(models.RegisterUserRequest)
	err = userContext.Bind(request)
	if err != nil {
		return userContext.JSON(http.StatusBadRequest, err.Error())
	}

	//client validation
	errs := validateUser(*request)
	if len(errs) > 0 {
		return userContext.JSON(http.StatusBadRequest, errs)
	}

	//generate api key
	userid := sequentialguid.NewSequentialGuid().String()
	usercreatedOn := time.Now()

	tx, _ := service.Db.Begin()

	//save user
	_, err = tx.Exec("INSERT INTO users VALUES(?,?,?,?,?,?);", userid, request.UserName, request.Email, usercreatedOn, usercreatedOn, usercreatedOn)
	if err != nil {
		return userContext.JSON(http.StatusInternalServerError, err.Error())
	}

	//save keys
	userKeyId := sequentialguid.NewSequentialGuid().String()
	apikey := helpers.GenerateApiKey(request.Email)
	keyCreatedOn := time.Now()
	expiryDate := keyCreatedOn.AddDate(expiryYear, 0, 0)

	_, err = tx.Exec("INSERT INTO userkeys VALUES(?,?,?,?,?,?);", userKeyId, apikey, keyCreatedOn, keyCreatedOn, expiryDate, userid, true)
	if err != nil {
		tx.Rollback()
		return userContext.JSON(http.StatusInternalServerError, err.Error())
	}

	tx.Commit()

	return userContext.JSON(http.StatusOK, models.RegisterUserResponse{Id: userid, ApiKey: userKeyId})
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