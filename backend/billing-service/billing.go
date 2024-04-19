package billing

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/samber/lo"
)

const (
	characterLimit = 8
	expirySpan     = 1
	dbFilePath     = "urlshortnerDB.db"
	RootFolderPath = "backend"
)

type BillingService struct {
	db            *sql.DB
	scheduleTimer *time.Ticker
}

type OrganizationInfo struct {
	Id       string
	Name     string
	Timezone string
}

type OrganizationPayPlanInfo struct {
	Id        string
	PlayType  enums.PayPlan
	PayPlanId string
	Amount    float64
	PayCycle  enums.PayCycle
	Status    enums.PlanStatus
}

type PayScheduleInfo struct {
	Id            string
	EffectiveDate time.Time
	EndDate       time.Time
}

func (service *BillingService) Run(done chan os.Signal) {
	fmt.Println("Starting the billing service..")
	for {
		select {
		case <-done:
			service.scheduleTimer.Stop()
			fmt.Println("Exiting billing service ....")
			return
		case t := <-service.scheduleTimer.C:
			fmt.Println("Tick at ....", t)
			service.generateBilling()
		}
	}
}

func (service *BillingService) generateBilling() {
	rows, err := service.db.Query("SELECT Id,Name,TimeZone FROM Organizations WHERE IsDeprecated =?", false)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var organizationInfos []OrganizationInfo
	for rows.Next() {
		var organizationInfo OrganizationInfo
		err = rows.Scan(&organizationInfo.Id, &organizationInfo.Name, &organizationInfo.Timezone)
		if err != nil {
			log.Fatal(err)
		}
		organizationInfos = append(organizationInfos, organizationInfo)
	}

	for _, organizationInfo := range organizationInfos {
		rows, err := service.db.Query(`SELECT organizationpayplans.Id, organizationpayplans.PayPlanId,organizationpayplans.PayCycle,payplans.Type,	
		 payplans.Amount,organizationpayplans.Status FROM payplans
			JOIN organizationpayplans ON payplans.Id = organizationpayplans.PayPlanId
			WHERE organizationpayplans.OrganizationId =?
			AND payplans.IsLatest =?
			AND (organizationpayplans.Status =? OR organizationpayplans.Status =?)`,
			organizationInfo.Id, true, enums.Current, enums.Upcoming)

		if err != nil {
			fmt.Printf("Error getting organization pay plan: %s for specific organization: %s \n", err.Error(), organizationInfo.Name)
			continue
		}

		var organizationPayPlanInfos []OrganizationPayPlanInfo
		for rows.Next() {
			var organizationPayPlanInfo OrganizationPayPlanInfo
			err = rows.Scan(&organizationPayPlanInfo.Id, &organizationPayPlanInfo.PayPlanId, &organizationPayPlanInfo.PayCycle,
				&organizationPayPlanInfo.PlayType, &organizationPayPlanInfo.Amount, &organizationPayPlanInfo.Status)
			if err != nil {
				fmt.Printf("Error scanning organization pay plan: %s for specific organization: %s \n", err.Error(), organizationInfo.Name)
				continue
			}
			organizationPayPlanInfos = append(organizationPayPlanInfos, organizationPayPlanInfo)
		}

		location, err := time.LoadLocation(organizationInfo.Timezone)
		if err != nil {
			fmt.Printf("Error loading location: %s for specific organization : %s \n", err.Error(), organizationInfo.Name)
			continue
		}

		currentTime := time.Now().In(location)
		fmt.Println("current time :", currentTime)

		currentOrganizationPlan, ok := lo.Find(organizationPayPlanInfos, func(info OrganizationPayPlanInfo) bool {
			return info.Status == enums.Current
		})

		if ok && !reflect.ValueOf(currentOrganizationPlan).IsZero() {
			if currentOrganizationPlan.PlayType == enums.Free {
				continue
			}

			var currentPayScheduleInfo PayScheduleInfo
			payScheduleRow := service.db.QueryRow("SELECT Id,EffectiveDate, EndDate FROM payschedules WHERE OrganizationId =? ORDER BY EffectiveDate DESC LIMIT 1", organizationInfo.Id)
			err = payScheduleRow.Scan(&currentPayScheduleInfo.Id, &currentPayScheduleInfo.EffectiveDate, &currentPayScheduleInfo.EndDate)
			if err != nil {
				fmt.Printf("Error getting organization pay schedules : %s for specific organization: %s \n", err.Error(), organizationInfo.Name)
				continue
			}

			tx, err := service.db.Begin()
			if err != nil {
				fmt.Println(err)
				continue
			}

			if !currentTime.After(currentPayScheduleInfo.EndDate) {
				continue
			}

			revenueId := sequentialguid.NewSequentialGuid().String()
			createdOn := time.Now()
			_, err = tx.Exec("INSERT INTO revenues VALUES(?,?,?,?,?,?,?,?);", revenueId, currentOrganizationPlan.Amount,
				currentPayScheduleInfo.EffectiveDate, currentPayScheduleInfo.EndDate, currentPayScheduleInfo.Id, organizationInfo.Id, createdOn, createdOn)
			if err != nil {
				tx.Rollback()
				fmt.Printf("Error insert revenue changes : %s for specific organization: %s \n", err.Error(), organizationInfo.Name)
			}

			upcomingOrganizationPlan, ok2 := lo.Find(organizationPayPlanInfos, func(info OrganizationPayPlanInfo) bool {
				return info.Status == enums.Upcoming
			})

			var newOrganizationPlanId string
			effectiveDate := currentPayScheduleInfo.EndDate.Add(1 * time.Second)
			var endDate time.Time

			if !ok2 && reflect.ValueOf(upcomingOrganizationPlan).IsZero() {
				newOrganizationPlanId = currentOrganizationPlan.Id
				if currentOrganizationPlan.PayCycle == enums.Monthly {
					endDate = effectiveDate.AddDate(0, 1, 0)
				} else if currentOrganizationPlan.PayCycle == enums.Yearly {
					endDate = effectiveDate.AddDate(1, 0, 0)
				}
			} else {
				if upcomingOrganizationPlan.PlayType == enums.Team {
					continue
				}

				if upcomingOrganizationPlan.PayCycle == enums.Monthly {
					endDate = effectiveDate.AddDate(0, 1, 0)
				} else if upcomingOrganizationPlan.PayCycle == enums.Yearly {
					endDate = effectiveDate.AddDate(1, 0, 0)
				}

				newOrganizationPlanId = upcomingOrganizationPlan.Id
				_, err = tx.Exec("UPDATE organizationpayplans SET status =? WHERE Id =?", enums.Archived, currentOrganizationPlan.Id)
				if err != nil {
					fmt.Println(err)
				}

				_, err = tx.Exec("UPDATE organizationpayplans SET status =? WHERE Id =?", enums.Current, upcomingOrganizationPlan.Id)
				if err != nil {
					fmt.Println(err)
				}
			}

			newPayScheduleId := sequentialguid.NewSequentialGuid().String()
			nextPayScheduleCreatedOn := time.Now()
			_, err = tx.Exec("INSERT INTO payschedules VALUES(?,?,?,?,?,?,?,?);", newPayScheduleId, effectiveDate,
				endDate, nextPayScheduleCreatedOn, nextPayScheduleCreatedOn, organizationInfo.Id, newOrganizationPlanId, false)
			if err != nil {
				tx.Rollback()
				fmt.Printf("Error insert revenue changes : %s for specific organization: %s \n", err.Error(), organizationInfo.Name)
			}

			tx.Commit()
		}
	}
}

func NewBillingService() (*BillingService, error) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	index := strings.Index(currentWorkingDirectory, RootFolderPath)
	if index == -1 {
		return nil, errors.New("app Root Folder Path not found")
	}

	filePath := filepath.Join(currentWorkingDirectory[:index], RootFolderPath, dbFilePath)
	db, err := database.ConnectToSQLite(filePath)
	if err != nil {
		return nil, err
	}

	return &BillingService{
		db:            db,
		scheduleTimer: time.NewTicker(24 * time.Hour),
	}, nil
}
