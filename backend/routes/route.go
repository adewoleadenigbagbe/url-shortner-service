package routes

import (
	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	middlewares "github.com/adewoleadenigbagbe/url-shortner-service/middleware"
)

func RegisterRoutes(app *core.BaseApp) {
	router := app.Echo
	router.POST("/api/v1/auth/register", app.AuthService.RegisterUser)
	router.POST("/api/v1/auth/sign-in", app.AuthService.LoginUser)
	router.POST("/api/v1/url", app.UrlService.CreateShortUrl, middlewares.AuthorizeUser)
}
