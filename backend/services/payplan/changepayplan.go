package services

import (
	"database/sql"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
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

	row := service.Db.QueryRow(`SELECT payplans.Id from payplans organizationpayplans 
	JOIN organizationpayplans ON payplans.Id = organizationpayplans.PayPlanId
	WHERE organizationpayplans.OrganizationId =? 
	AND organizationpayplans.PayPlanId =? AND organizationpayplans.IsLatest =?
	AND payplans.IsLatest =?`,
		request.OrganizationId, request.PayplanId, true, true, true)

	var planId string
	err = row.Scan(&planId)
	if err != nil {
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if planId == request.PayplanId {
		return planContext.JSON(http.StatusBadRequest, []string{"plan already exist."})
	}

	_, err = tx.Exec("UPDATE organizationpayplans SET IsLatest =? WHERE OrganizationId =?", false, request.OrganizationId)
	if err != nil {
		tx.Rollback()
	}

	id := sequentialguid.NewSequentialGuid().String()
	now := time.Now()
	_, err = tx.Exec("INSERT INTO organizationpayplans VALUES(?,?,?,?,?,?,?);",
		id, request.PayCycle, request.PayplanId, request.OrganizationId, now, now, true)

	if err != nil {
		tx.Rollback()
		return planContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	tx.Commit()
	return nil
}
