package services

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/labstack/echo/v4"
)

func (service EnumService) GetPayCycle(enumContext echo.Context) error {
	var d enums.PayCycle

	values := d.GetKeyValues()
	return enumContext.JSON(http.StatusOK, values)
}
