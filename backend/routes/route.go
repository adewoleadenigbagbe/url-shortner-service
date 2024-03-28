package routes

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	middlewares "github.com/adewoleadenigbagbe/url-shortner-service/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(app *core.BaseApp, middleware *middlewares.AppMiddleware) {
	router := app.Echo
	router.POST("/api/v1/auth/register", app.AuthService.RegisterUser)
	router.POST("/api/v1/auth/sign-in", app.AuthService.LoginUser)
	router.POST("/api/v1/auth/sign-out", app.AuthService.LogOut)

	router.POST("/api/v1/shortlink", app.UrlService.CreateShortUrl, middleware.AuthorizeUser)
	router.GET("/api/v1/shortlink", app.UrlService.GetShortLinks, middleware.AuthorizeUser)
	router.GET("/api/v1/shortlink/redirect", app.UrlService.RedirectShort)

	router.POST("/api/v1/domain", app.DomainService.CreateDomain, middleware.AuthorizeAdmin)

	router.POST("/api/v1/user/send-email", app.UserService.SendEmail, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission, middleware.AuthorizeFeaturePermission)
	router.POST("/api/v1/user/convert-referral", app.UserService.ConvertReferral)

	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
