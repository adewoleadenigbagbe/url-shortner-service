package services

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type TagService struct {
	Db *sql.DB
}

type tagInfo struct {
	Id   string
	Name string
}

func (service TagService) CreateTag(tagContext echo.Context) error {
	var (
		err error
	)

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

	request.Tags = lo.Map(request.Tags, func(tag string, index int) string {
		return strings.ToLower(tag)
	})

	queryStr, args := formatTagQuery(request.Tags)
	rows, err := service.Db.Query(queryStr, args...)
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	defer rows.Close()

	var existingTagInfos []tagInfo
	for rows.Next() {
		var tag tagInfo
		err = rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		existingTagInfos = append(existingTagInfos, tag)
	}

	if len(existingTagInfos) > 0 {
		existingTags := lo.Map(existingTagInfos, func(tag tagInfo, index int) string {
			return tag.Name
		})
		request.Tags, _ = lo.Difference(request.Tags, existingTags)
	}

	if len(request.Tags) > 0 {
		tagInfos := lo.Map(request.Tags, func(tag string, index int) tagInfo {
			return tagInfo{
				Id:   sequentialguid.NewSequentialGuid().String(),
				Name: tag,
			}
		})

		stmt, args := insertTagStmt(tagInfos)
		_, err = service.Db.Exec(stmt, args...)
		if err != nil {
			return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
	}

	return tagContext.JSON(http.StatusCreated, nil)
}

func validateTagRequest(request models.CreateTagRequest) []error {
	var validationErrors []error

	if len(request.Tags) == 0 {
		validationErrors = append(validationErrors, errors.New("tags is required. supply at least one"))
	}
	return validationErrors
}

func formatTagQuery(tags []string) (string, []interface{}) {
	str := "SELECT Id, Name FROM tags WHERE Name IN (?" + strings.Repeat(",?", len(tags)-1) + ")"
	vals := []interface{}{}
	for _, tag := range tags {
		vals = append(vals, tag)
	}
	return str, vals
}

func insertTagStmt(tags []tagInfo) (string, []interface{}) {
	stmt := "INSERT INTO tags VALUES "
	vals := []interface{}{}
	for _, tag := range tags {
		stmt += "(?,?,?),"
		now := time.Now()
		vals = append(vals, tag.Id, tag.Name, now)
	}
	//trim the last ,
	stmt = stmt[0 : len(stmt)-1]
	return stmt, vals
}
