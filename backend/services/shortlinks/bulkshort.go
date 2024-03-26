package services

import (
	"fmt"
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/labstack/echo/v4"
)

func (service UrlService) CreateBulkShortLink(urlContext echo.Context) error {
	var err error
	req := urlContext.Request()
	contentType := req.Header.Get("Content-Type")
	body := req.Body
	reader, err := helpers.CreateReader(contentType, body)
	if err != nil {
		urlContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	data, err := reader.ReadFile()
	if err != nil {
		urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	fmt.Println(data)

	return nil
}
