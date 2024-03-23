package services

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

func (service DomainService) DeleteDomain(domainContext echo.Context) error {
	var err error
	request := new(models.DeleteDomainRequest)
	err = domainContext.Bind(request)
	if err != nil {
		return domainContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	if len(request.Name) == 0 {
		return domainContext.JSON(http.StatusBadRequest, []string{"domain name is required"})
	}

	result, _ := service.Db.Exec("UPDATE domains SET IsDeprecated =? WHERE Name =?", false, request.Name)
	rows, err := result.RowsAffected()
	if err == nil && rows > 0 {
		return domainContext.JSON(http.StatusNoContent, nil)
	}

	return domainContext.JSON(http.StatusNotFound, []string{err.Error()})
}
