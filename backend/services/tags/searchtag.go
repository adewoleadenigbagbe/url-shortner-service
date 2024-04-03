package services

import (
	"fmt"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service TagService) SearchTag(tagContext echo.Context) error {
	var (
		err error
	)

	request := new(models.SearchTagRequest)
	err = tagContext.Bind(request)
	if err != nil {
		return tagContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	term := "'%" + request.Term + "%'"
	query := fmt.Sprintf("SELECT Name FROM tags WHERE Name LIKE %s", term)
	rows, err := service.Db.Query(query)
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		tags = append(tags, tag)
	}

	return tagContext.JSON(http.StatusOK, tags)
}
