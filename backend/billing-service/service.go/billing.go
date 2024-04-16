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
}

type PayScheduleInfo struct {
	Id            string
	EffectiveDate time.Time
	EndDate       time.Time
}

func (service *BillingService) Run(done chan os.Signal) {
	fmt.Println("Starting the key generation service..")
	for {
		select {
		case <-done:
			service.scheduleTimer.Stop()
			fmt.Println("Exiting key generator service ....")
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
		err = rows.Scan(&organizationInfo.Id, organizationInfo.Name, &organizationInfo.Timezone)
		if err != nil {
			log.Fatal(err)
		}
		organizationInfos = append(organizationInfos, organizationInfo)
	}

	for _, organizationInfo := range organizationInfos {
		row := service.db.QueryRow(`SELECT organizationpayplans.Id, organizationpayplans.PayPlanId,organizationpayplans.PayCycle,payplans.Type, payplans.Amount FROM payplans 
			JOIN organizationpayplans ON payplans.Id = organizationpayplans.PayPlanId
			WHERE organizationpayplans.OrganizationId =? 
			AND organizationpayplans.Status =? 
			AND payplans.IsLatest =?`,
			organizationInfo.Id, enums.Current, true)

		var organizationPayPlanInfo OrganizationPayPlanInfo
		err = row.Scan(&organizationPayPlanInfo.Id, &organizationPayPlanInfo.PayPlanId, &organizationPayPlanInfo.PayCycle, &organizationPayPlanInfo.PlayType, &organizationPayPlanInfo.Amount)
		if err != nil {
			fmt.Printf("Error getting organization pay plan: %s for specific organization: %s ", err.Error(), organizationInfo.Name)
			continue
		}

		if organizationPayPlanInfo.PlayType == enums.Free {
			continue
		}

		location, err := time.LoadLocation(organizationInfo.Timezone)
		if err != nil {
			fmt.Printf("Error loading location: %s for specific organization : %s ", err.Error(), organizationInfo.Name)
			continue
		}

		currentTime := time.Now().In(location)
		fmt.Println(currentTime)

		var payScheduleInfo PayScheduleInfo
		row2 := service.db.QueryRow("SELECT Id,EffectiveDate, EndDate FROM payschedules WHERE OrganizationId =? AND IsNext =? ORDER BY EffectiveDate LIMIT 1", organizationInfo.Id, false)
		err = row2.Scan(&payScheduleInfo.Id, &payScheduleInfo.EffectiveDate, &payScheduleInfo.EndDate)
		if err != nil {
			fmt.Printf("Error getting organization pay schedules : %s for specific organization: %s ", err.Error(), organizationInfo.Name)
			continue
		}

		var nextPayScheduleInfo PayScheduleInfo
		row3 := service.db.QueryRow("SELECT Id,EffectiveDate, EndDate FROM payschedules WHERE OrganizationId =? AND IsNext =? ORDER BY EffectiveDate LIMIT 1", organizationInfo.Id, true)
		err = row3.Scan(&nextPayScheduleInfo.Id, &nextPayScheduleInfo.EffectiveDate, &nextPayScheduleInfo.EndDate)
		if err != nil {
			fmt.Printf("Error getting organization pay schedules : %s for specific organization: %s ", err.Error(), organizationInfo.Name)
			continue
		}

		tx, err := service.db.Begin()
		if err != nil {
			fmt.Println(err)
		}

		if reflect.ValueOf(nextPayScheduleInfo).IsZero() {
			newPayScheduleId := sequentialguid.NewSequentialGuid().String()
			effectiveDate := payScheduleInfo.EndDate.Add(1 * time.Second)
			var endDate time.Time
			nextPayScheduleCreatedOn := time.Now()

			if organizationPayPlanInfo.PayCycle == enums.Monthly {
				endDate = effectiveDate.AddDate(0, 1, 0)
			} else if organizationPayPlanInfo.PayCycle == enums.Yearly {
				endDate = effectiveDate.AddDate(1, 0, 0)
			}

			_, err = tx.Exec("INSERT INTO payschedules VALUES(?,?,?,?,?,?,?,?);", newPayScheduleId, effectiveDate,
				endDate, nextPayScheduleCreatedOn, nextPayScheduleCreatedOn, organizationInfo.Id, organizationPayPlanInfo.Id, false)
			if err != nil {
				tx.Rollback()
				fmt.Printf("Error insert revenue changes : %s for specific organization: %s ", err.Error(), organizationInfo.Name)
			}
		}

		if currentTime.After(payScheduleInfo.EndDate) {
			revenueId := sequentialguid.NewSequentialGuid().String()
			createdOn := time.Now()
			_, err = tx.Exec("INSERT INTO revenues VALUES(?,?,?,?,?,?,?,?);", revenueId, organizationPayPlanInfo.Amount,
				payScheduleInfo.EffectiveDate, payScheduleInfo.EndDate, payScheduleInfo.Id, organizationInfo.Id, createdOn, createdOn)
			if err != nil {
				tx.Rollback()
				fmt.Printf("Error insert revenue changes : %s for specific organization: %s ", err.Error(), organizationInfo.Name)
			}
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
