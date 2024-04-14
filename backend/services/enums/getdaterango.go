package services

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/labstack/echo/v4"
)

type EnumService struct {
}

func (service EnumService) GetDateRanges(enumContext echo.Context) error {
	var d enums.DateRange

	values := d.GetValues()
	return enumContext.JSON(http.StatusOK, values)
}
