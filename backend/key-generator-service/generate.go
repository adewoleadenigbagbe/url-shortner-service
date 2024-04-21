package linkservice

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
)

const (
	characterLimit = 8
	expirySpan     = 1
	dbFilePath     = "data/urlshortnerDB.db"
	RootFolderPath = "backend"
)

type ShortlinkGenerator struct {
	db            *sql.DB
	scheduleTimer *time.Ticker
	done          chan bool
}

func (sg *ShortlinkGenerator) GenerateLink(done chan os.Signal) {
	fmt.Println("Starting the key generation service..")
	for {
		select {
		case <-done:
			sg.scheduleTimer.Stop()
			fmt.Println("Exiting key generator service ....")
			return
		case t := <-sg.scheduleTimer.C:
			fmt.Println("Tick at ....", t)
			shortKey := helpers.RandStringBytesRmndr(characterLimit)
			err := sg.insertlink(shortKey)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (sg *ShortlinkGenerator) insertlink(key string) error {
	now := time.Now()
	expirationDate := now.AddDate(expirySpan, 0, 0)
	_, err := sg.db.Exec("INSERT INTO unusedshortlinks VALUES(?,?,?,?,?);",
		key, now, now, expirationDate, false)

	return err
}

func NewShortlinkGenerator() *ShortlinkGenerator {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	index := strings.Index(currentWorkingDirectory, RootFolderPath)
	if index == -1 {
		log.Fatal("App Root Folder Path not found")
	}

	filePath := filepath.Join(currentWorkingDirectory[:index], RootFolderPath, dbFilePath)
	db, err := database.ConnectToSQLite(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return &ShortlinkGenerator{
		db:            db,
		scheduleTimer: time.NewTicker(2 * time.Minute),
		done:          make(chan bool),
	}
}
