package services

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

const (
	EmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
)

func RegisterUser(userContext echo.Context) error {
	var err error
	request := new(models.UserRequest)
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

	//save user

	return nil
}

func validateUser(user models.UserRequest) []error {
	var validationErrors []error

	if user.UserName == "" {
		validationErrors = append(validationErrors, fmt.Errorf("username is required"))
	}

	isEmailValid, _ := regexp.MatchString(EmailRegex, user.Email)
	if !isEmailValid {
		validationErrors = append(validationErrors, fmt.Errorf("email is invalid"))
	}
	return validationErrors
}
