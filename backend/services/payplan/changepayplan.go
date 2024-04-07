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
	tx, err := service.Db.Begin()
	if err != nil {
		return err
	}

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

	row := service.Db.QueryRow(`SELECT organizations.TimeZone payplans.Id,payplans.Type from organizations
	JOIN organizationpayplans ON organizations.Id = organizationpayplans.OrganizationId
	JOIN payplans ON organizationpayplans.PayPlanId = payplans.Id
	WHERE organizations.Id =? 
	AND organizationpayplans.IsLatest =?
	AND payplans.IsLatest =?`,
		request.OrganizationId, true, true)

	var timezone string
	var planId string
	var planType enums.PayPlan
	err = row.Scan(&timezone, &planId, &planType)
	if err != nil {
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if planId == request.PayplanId {
		return planContext.JSON(http.StatusBadRequest, []string{"plan already exist."})
	}

	_, err = tx.Exec("UPDATE organizationpayplans SET IsLatest =? WHERE OrganizationId =?", false, request.OrganizationId)
	if err != nil {
		tx.Rollback()
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if planType == enums.Free {
		request.PayCycle = enums.None
	}

	organizationpayplanId := sequentialguid.NewSequentialGuid().String()
	now := time.Now()
	_, err = tx.Exec("INSERT INTO organizationpayplans VALUES(?,?,?,?,?,?,?);",
		organizationpayplanId, request.PayCycle, request.PayplanId, request.OrganizationId, now, now, true)
	if err != nil {
		tx.Rollback()
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if planType != enums.Free {
		location, err := time.LoadLocation(timezone)
		if err != nil {
			return planContext.JSON(http.StatusBadRequest, []string{err.Error()})
		}

		effectiveDate := time.Now().In(location)
		var endDate time.Time
		var nextEndDate time.Time
		if request.PayCycle == enums.Monthly {
			endDate = effectiveDate.AddDate(0, 1, 0)
			nextEndDate = endDate.AddDate(0, 1, 0)
		} else if request.PayCycle == enums.Yearly {
			endDate = effectiveDate.AddDate(1, 0, 0)
			nextEndDate = endDate.AddDate(1, 0, 0)
		}

		payScheduleId := sequentialguid.NewSequentialGuid().String()
		payScheduleCreatedOn := time.Now()
		_, err = tx.Exec("INSERT INTO payschedules VALUES(?,?,?,?,?,?,?,?);",
			payScheduleId, effectiveDate, endDate, payScheduleCreatedOn, payScheduleCreatedOn, request.OrganizationId, organizationpayplanId, false)
		if err != nil {
			tx.Rollback()
			return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}

		nextPayScheduleId := sequentialguid.NewSequentialGuid().String()
		nextPayScheduleCreatedOn := time.Now()
		_, err = tx.Exec("INSERT INTO payschedules VALUES(?,?,?,?,?,?,?,?);",
			nextPayScheduleId, endDate.Add(1*time.Second), nextEndDate, nextPayScheduleCreatedOn, nextPayScheduleCreatedOn, request.OrganizationId, organizationpayplanId, true)
		if err != nil {
			tx.Rollback()
			return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}

	}

	tx.Commit()
	return nil
}
