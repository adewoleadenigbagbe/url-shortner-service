package services

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/labstack/echo/v4"
)

func (service EnumService) GetPayPlan(enumContext echo.Context) error {
	var d enums.PayPlan

	values := d.GetValues()
	return enumContext.JSON(http.StatusOK, values)
}
