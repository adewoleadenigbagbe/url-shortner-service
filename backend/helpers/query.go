package helpers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func IsCurrentUser(db *sql.DB, id, key string) bool {
	query := `SELECT COUNT(1) FROM users JOIN userkeys on users.Id = userkeys.UserId 
	WHERE users.Id =? AND userkeys.ApiKey =? AND userkeys.IsActive =? AND userkeys.ExpirationDate >?`
	var count int64
	err := db.QueryRow(query, id, key, true, time.Now()).Scan(&count)
	if err != nil {
		return false
	}
	return count == 1
}

func CheckForBlackListedTokens(context echo.Context, redisClient *redis.Client, id string) bool {
	token := GetTokenFromRequest(context)
	res, err := redisClient.Get(context.Request().Context(), id).Result()
	if err != redis.Nil {
		fmt.Println(err)
	}
	return res == token
}

func CheckUserOrganization(db *sql.DB, id string) (bool, error) {
	query := `SELECT COUNT(1) FROM organizations JOIN users ON organizations.Id = users.OrganizationId 
	WHERE organizations.Id=? AND organizations.IsDeprecated=? AND users.IsDeprecated=?`
	var count int64
	err := db.QueryRow(query, id, false, false).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckPayPlan(db *sql.DB, id string) (enums.PayPlan, error) {
	query := `SELECT payplans.Type FROM payplans JOIN organizationpayplans ON  payplans.Id = organizationpayplans.PayPlanId
	WHERE organizationpayplans.OrganizationId=? AND payplans.IsLatest=? AND organizationpayplans.IsLatest=?`
	var planType enums.PayPlan
	err := db.QueryRow(query, id, true, true).Scan(&planType)
	if err != nil {
		return 0, err
	}
	return planType, nil
}
