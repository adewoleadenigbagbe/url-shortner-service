package core

import jwtauth "github.com/adewoleadenigbagbe/url-shortner-service/helpers/auth"

func RegisterRoutes(app *BaseApp) {
	router := app.echo
	router.POST("/api/v1/auth/register", app.AuthService.RegisterUser)
	router.POST("/api/v1/auth/sign-in", app.AuthService.LoginUser)
	router.POST("/api/v1/url", app.UrlService.CreateShortUrl, jwtauth.AuthorizeUser)
}
