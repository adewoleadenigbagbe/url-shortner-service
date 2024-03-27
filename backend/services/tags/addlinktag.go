package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func (service TagService) AddShortLinkTag(tagContext echo.Context) error {
	var err error
	request := new(models.CreateShortLinkTagRequest)
	err = tagContext.Bind(request)
	if err != nil {
		return tagContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	modelErrs := validateShortLinkTagRequest(*request)
	if len(modelErrs) > 0 {
		valErrors := lo.Map(modelErrs, func(er error, index int) string {
			return er.Error()
		})
		return tagContext.JSON(http.StatusBadRequest, valErrors)
	}

	var count int
	row := service.Db.QueryRow("SELECT COUNT(1) FROM shortlinktags WHERE TagId =?", request.TagId)
	err = row.Scan(&count)
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	fmt.Println("count :", count)
	if count > 0 {
		return tagContext.JSON(http.StatusBadRequest, []string{"tag name already exist"})
	}

	id := sequentialguid.NewSequentialGuid().String()
	now := time.Now()
	_, err = service.Db.Exec("INSERT INTO shortlinktags VALUES(?,?,?,?);", id, request.Short, request.TagId, now)
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	return tagContext.JSON(http.StatusCreated, models.CreateShortLinkTagResponse{Id: id})
}

func validateShortLinkTagRequest(request models.CreateShortLinkTagRequest) []error {
	var validationErrors []error

	if request.Short == "" {
		validationErrors = append(validationErrors, errors.New("shortlink is required"))
	}

	if request.TagId == "" {
		validationErrors = append(validationErrors, errors.New("tagId is required"))
	}
	return validationErrors
}
