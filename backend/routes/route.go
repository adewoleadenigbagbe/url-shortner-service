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

	//Shortlinks
	router.POST("/api/v1/shortlink", app.UrlService.CreateShortLink, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission)
	router.GET("/api/v1/shortlink", app.UrlService.GetShortLinks, middleware.AuthorizeUser)
	router.POST("/api/v1/shortlink/redirect", app.UrlService.RedirectShort)
	router.POST("/api/v1/shortlink/bulk", app.UrlService.CreateBulkShortLink)

	//Domains
	router.POST("/api/v1/domain", app.DomainService.CreateDomain, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission)
	router.GET("/api/v1/domain", app.DomainService.GetDomains)

	//Users
	router.POST("/api/v1/user/send-email", app.UserService.SendEmail, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission, middleware.AuthorizeFeaturePermission)
	router.POST("/api/v1/user/convert-referral", app.UserService.ConvertReferral)

	//Teams
	router.POST("/api/v1/teams", app.TeamService.AddTeam, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission, middleware.AuthorizeFeaturePermission)
	router.POST("/api/v1/teams/add-user", app.TeamService.AddUserToTeam, middleware.AuthorizeAdmin)
	router.GET("/api/v1/teams/search", app.TeamService.SearchTeam, middleware.AuthorizeAdmin, middleware.AuthourizeOrganizationPermission, middleware.AuthorizeFeaturePermission)

	//Tags
	router.POST("/api/v1/tags", app.TagService.CreateTag, middleware.AuthorizeAdmin)
	router.POST("/api/v1/tags/add-tag-short", app.TagService.AddShortLinkTag, middleware.AuthorizeAdmin)
	router.GET("/api/v1/tags/search", app.TagService.SearchTag)

	//Report
	router.POST("/api/v1/report/links", app.ExportService.GenerateShortLinkReport)

	//statistics
	router.GET("/api/v1/statistics", app.StatisticsService.GetShortStatistics)

	//enums
	router.GET("/api/v1/enum/date-range", app.EnumService.GetDateRanges)
	router.GET("/api/v1/enum/pay-cycle", app.EnumService.GetPayCycle)
	router.GET("/api/v1/enum/pay-plan", app.EnumService.GetPayPlan)
	router.GET("/api/v1/enum/report-type", app.EnumService.GetReportType)
	router.GET("/api/v1/enum/role", app.EnumService.GetRoles)

	//default path
	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
