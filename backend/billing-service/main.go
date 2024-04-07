package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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

type PayPlanInfo struct {
	PlayType  enums.PayPlan
	PayPlanId string
}

func (service *BillingService) Run() {
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
		row := service.db.QueryRow(`SELECT organizationpayplans.PayPlanId,payplans.Type FROM payplans 
			JOIN organizationpayplans ON payplans.Id = organizationpayplans.PayPlanId
			WHERE organizationpayplans.OrganizationId =? 
			AND organizationpayplans.IsLatest =? 
			AND payplans.IsLatest =?`,
			organizationInfo.Id, true, true)

		var payPlanInfo PayPlanInfo
		err = row.Scan(&payPlanInfo.PlayType, &payPlanInfo.PayPlanId)
		if err != nil {
			fmt.Printf("Error getting organization pay plan: %s for specific organization: %s ", err.Error(), organizationInfo.Name)
			continue
		}

		if payPlanInfo.PlayType == enums.Free {
			continue
		}

		location, err := time.LoadLocation(organizationInfo.Timezone)
		if err != nil {
			fmt.Printf("Error loading location: %s for specific organization: %s ", err.Error(), organizationInfo.Name)
			continue
		}

		currentTime := time.Now().In(location)
		fmt.Println(currentTime)
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

func main() {
	service, err := NewBillingService()
	if err != nil {
		log.Fatal(err)
	}

	service.Run()
}
