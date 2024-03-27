package services

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type TagService struct {
	Db *sql.DB
}

func (service TagService) CreateTag(tagContext echo.Context) error {
	var err error
	request := new(models.CreateTagRequest)
	err = tagContext.Bind(request)
	if err != nil {
		return tagContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	modelErrs := validateTagRequest(*request)
	if len(modelErrs) > 0 {
		valErrors := lo.Map(modelErrs, func(er error, index int) string {
			return er.Error()
		})
		return tagContext.JSON(http.StatusBadRequest, valErrors)
	}

	var count int
	row := service.Db.QueryRow("SELECT COUNT(1) FROM tags WHERE tags =?", request.Name)
	err = row.Scan(&count)
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	fmt.Println("count :", count)
	if count > 0 {
		return tagContext.JSON(http.StatusBadRequest, []string{"tag name already exist"})
	}

	tx, err := service.Db.Begin()
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	tagId := sequentialguid.NewSequentialGuid().String()
	tagCreatedOn := time.Now()
	_, err = service.Db.Exec("INSERT INTO tags VALUES(?,?,?);", tagId, request.Name, tagCreatedOn)
	if err != nil {
		tx.Rollback()
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	shortlinktagId := sequentialguid.NewSequentialGuid().String()
	linktagCreatedOn := time.Now()
	_, err = service.Db.Exec("INSERT INTO shortlinktags VALUES(?,?,?,?);", shortlinktagId, request.Short, tagId, linktagCreatedOn)
	if err != nil {
		tx.Rollback()
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	tx.Commit()
	return tagContext.JSON(http.StatusCreated, models.CreateTagResponse{Id: tagId})
}

func validateTagRequest(request models.CreateTagRequest) []error {
	var validationErrors []error
	if request.Short == "" {
		validationErrors = append(validationErrors, errors.New("shortlink is required"))
	}

	if request.Name == "" {
		validationErrors = append(validationErrors, errors.New("name is required"))
	}
	return validationErrors
}
