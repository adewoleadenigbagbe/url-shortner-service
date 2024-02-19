package main

import (
	services "github.com/adewoleadenigbagbe/url-shortner-service/service"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(route *echo.Echo) {
	route.POST("/api/v1/auth/register", services.RegisterUser)
	route.POST("/api/v1/auth/sign-in", services.LoginUser)
}
