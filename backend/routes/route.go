package routes

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	middlewares "github.com/adewoleadenigbagbe/url-shortner-service/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(app *core.BaseApp, middleware *middlewares.AppMiddleware) {
	router := app.Echo

	//auth
	router.POST("/api/v1/auth/register", app.AuthService.RegisterUser)
	router.POST("/api/v1/auth/sign-in", app.AuthService.LoginUser)
	router.POST("/api/v1/auth/sign-out", app.AuthService.LogOut)

	//shortlinks
	router.POST("/api/v1/shortlink", app.UrlService.CreateShortLink, middleware.AuthorizeUser)
	router.GET("/api/v1/shortlink", app.UrlService.GetShortLinks, middleware.AuthorizeUser)
	router.GET("/api/v1/shortlink/redirect", app.UrlService.RedirectShort)

	//domains
	router.POST("/api/v1/domain", app.DomainService.CreateDomain, middleware.AuthorizeAdmin)

	//tags
	router.POST("/api/v1/tags", app.TagService.CreateTag, middleware.AuthorizeAdmin)
	router.POST("/api/v1/tags/add-tag-short", app.TagService.AddShortLinkTag, middleware.AuthorizeAdmin)

	//default path
	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
