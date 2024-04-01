package routes

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/core"
	middlewares "github.com/adewoleadenigbagbe/url-shortner-service/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(app *core.BaseApp, middleware *middlewares.AppMiddleware) {
	router := app.Echo

	//Auth
	router.POST("/api/v1/auth/register", app.AuthService.RegisterUser)
	router.POST("/api/v1/auth/sign-in", app.AuthService.LoginUser)
	router.POST("/api/v1/auth/sign-out", app.AuthService.LogOut)

	//shortlinks
	router.POST("/api/v1/shortlink", app.UrlService.CreateShortLink, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission)
	router.GET("/api/v1/shortlink", app.UrlService.GetShortLinks, middleware.AuthorizeUser)
	router.POST("/api/v1/shortlink/redirect", app.UrlService.RedirectShort)

	//domains
	router.POST("/api/v1/domain", app.DomainService.CreateDomain, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission)

	//users
	router.POST("/api/v1/user/send-email", app.UserService.SendEmail, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission, middleware.AuthorizeFeaturePermission)
	router.POST("/api/v1/user/convert-referral", app.UserService.ConvertReferral)

	//teams
	router.POST("/api/v1/teams", app.TeamService.AddTeam, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission, middleware.AuthorizeFeaturePermission)
	router.POST("/api/v1/teams/add-user", app.TeamService.AddUserToTeam, middleware.AuthorizeAdmin)
	router.GET("/api/v1/teams/search", app.TeamService.SearchTeam, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission, middleware.AuthorizeFeaturePermission)

	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
