package services

import (
	"database/sql"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
)

type PlanService struct {
	Db *sql.DB
}

func (service PlanService) ChangePayPlan(planContext echo.Context) error {
	var err error

	request := new(models.CreateOrganizationPlanRequest)
	binder := &echo.DefaultBinder{}
	err = binder.BindHeaders(planContext, request)
	if err != nil {
		return planContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	err = binder.BindBody(planContext, request)
	if err != nil {
		return planContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	row := service.Db.QueryRow(`SELECT organizations.TimeZone, payplans.Id,payplans.Type from organizations
	JOIN organizationpayplans ON organizations.Id = organizationpayplans.OrganizationId
	JOIN payplans ON organizationpayplans.PayPlanId = payplans.Id
	WHERE organizations.Id =? 
	AND organizationpayplans.Status =?
	AND payplans.IsLatest =?`,
		request.OrganizationId, enums.Current, true)

	var timezone string
	var existingPlanId string
	var existingPlanType enums.PayPlan
	err = row.Scan(&timezone, &existingPlanId, &existingPlanType)
	if err != nil {
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if existingPlanId == request.PayplanId {
		return planContext.JSON(http.StatusBadRequest, []string{"plan already exist."})
	}

	var planType enums.PayPlan
	row2 := service.Db.QueryRow("SELECT Type FROM payplans WHERE Id =? AND IsLatest =?", request.PayplanId, true)
	err = row2.Scan(&planType)
	if err != nil {
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	tx, err := service.Db.Begin()
	if err != nil {
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}
	_, err = tx.Exec("UPDATE organizationpayplans SET Status =? WHERE OrganizationId =? AND Status =?", enums.Archived, request.OrganizationId, enums.Upcoming)
	if err != nil {
		tx.Rollback()
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if planType == enums.Free {
		request.PayCycle = enums.None
	}

	organizationpayplanId := sequentialguid.NewSequentialGuid().String()
	var status enums.PlanStatus = enums.Upcoming
	if existingPlanType == enums.Free {
		_, err = tx.Exec("UPDATE organizationpayplans SET Status =? WHERE OrganizationId =? AND Status =?", enums.Archived, request.OrganizationId, enums.Current)
		if err != nil {
			tx.Rollback()
			return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}

		location, err := time.LoadLocation(timezone)
		if err != nil {
			return planContext.JSON(http.StatusBadRequest, []string{err.Error()})
		}

		currentTime := time.Now().In(location)
		var newEffectiveDate time.Time = currentTime
		var newEndDate time.Time

		if request.PayCycle == enums.Monthly {
			newEndDate = newEffectiveDate.AddDate(0, 1, 0)
		} else if request.PayCycle == enums.Yearly {
			newEndDate = newEffectiveDate.AddDate(1, 0, 0)
		}

		payScheduleId := sequentialguid.NewSequentialGuid().String()
		payScheduleCreatedOn := time.Now()
		_, err = tx.Exec("INSERT INTO payschedules VALUES(?,?,?,?,?,?,?,?);",
			payScheduleId, newEffectiveDate, newEndDate, payScheduleCreatedOn, payScheduleCreatedOn, request.OrganizationId, organizationpayplanId, false)
		if err != nil {
			tx.Rollback()
			return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}

		status = enums.Current
	}

	now := time.Now()
	_, err = tx.Exec("INSERT INTO organizationpayplans VALUES(?,?,?,?,?,?,?);",
		organizationpayplanId, request.PayCycle, request.PayplanId, request.OrganizationId, now, now, status)
	if err != nil {
		tx.Rollback()
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	tx.Commit()
	return planContext.JSON(http.StatusOK, models.CreateOrganizationPlanResponse{OrganizationPayPlanId: organizationpayplanId})
}
