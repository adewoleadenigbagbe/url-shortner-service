package services

import (
	"github.com/labstack/echo/v4"
)

func (service AuthService) LogOut(authContext echo.Context) error {
	//Invalidate user access token here
	return nil
}
