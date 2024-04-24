package middlewares

import (
	"net/http"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/labstack/echo/v4"
)

func (appMiddleware *AppMiddleware) AuthorizeUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		var err error
		id, err := helpers.ValidateUserJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "Authentication required")
		}

		tokenExist := helpers.CheckForBlackListedTokens(context, appMiddleware.Rdb, id)
		if tokenExist {
			return context.JSON(http.StatusBadRequest, "Invalid Authourization Token")
		}

		apikey := context.Request().Header.Get("X-Api-Key")
		if len(apikey) == 0 {
			return context.JSON(http.StatusBadRequest, "User ApiKey missing in the header")
		}

		if !helpers.IsCurrentUser(appMiddleware.Db, id, apikey) {
			return context.JSON(http.StatusUnauthorized, "invalid api key")
		}
		return next(context)
	}
}

func (appMiddleware *AppMiddleware) AuthorizeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		id, err := helpers.ValidateAdminRoleJWT(context)
		if err != nil {
			return context.JSON(http.StatusUnauthorized, "You are not allowed to access this resource")
		}

		tokenExist := helpers.CheckForBlackListedTokens(context, appMiddleware.Rdb, id)
		if tokenExist {
			return context.JSON(http.StatusBadRequest, "Invalid Authourization Token")
		}

		apikey := context.Request().Header.Get("X-Api-Key")
		if len(apikey) == 0 {
			return context.JSON(http.StatusBadRequest, "User ApiKey missing in the header")
		}

		if !helpers.IsCurrentUser(appMiddleware.Db, id, apikey) {
			return context.JSON(http.StatusUnauthorized, "invalid api key")
		}

		return next(context)
	}
}

func (appMiddleware *AppMiddleware) AuthourizeOrganizationPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		organizationId := context.Request().Header.Get("X-OrganizationId")
		exist, err := helpers.CheckUserOrganization(appMiddleware.Db, organizationId)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
		if !exist {
			return context.JSON(http.StatusUnauthorized, "no valid permission for this organization")
		}

		return next(context)
	}
}

func (appMiddleware *AppMiddleware) AuthorizeFeaturePermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		organizationId := context.Request().Header.Get("X-OrganizationId")
		planType, err := helpers.CheckPayPlan(appMiddleware.Db, organizationId)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, err.Error())
		}
		if planType == enums.Free {
			return context.JSON(http.StatusUnauthorized, "no valid permission for this paid feature")
		}

		return next(context)
	}
}
