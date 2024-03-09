package keyservice

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
	dbFilePath     = "urlshortnerDB.db"
	RootFolderPath = "backend"
)

type KeyGenerator struct {
	db            *sql.DB
	scheduleTimer *time.Ticker
	done          chan bool
}

func (kg *KeyGenerator) Run(done chan os.Signal) {
	for {
		select {
		case <-done:
			kg.scheduleTimer.Stop()
			fmt.Println("Exiting key generator service ....")
			return
		case t := <-kg.scheduleTimer.C:
			fmt.Println("Tick at ....", t)
			shortKey := helpers.RandStringBytesRmndr(characterLimit)
			err := kg.saveKey(shortKey)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (kg *KeyGenerator) saveKey(key string) error {
	now := time.Now()
	expirationDate := now.AddDate(expirySpan, 0, 0)
	_, err := kg.db.Exec("INSERT INTO unusedshortlinks VALUES(?,?,?,?,?);",
		key, now, now, expirationDate, false)

	return err
}

func NewKeyGenerator() *KeyGenerator {
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

	return &KeyGenerator{
		db:            db,
		scheduleTimer: time.NewTicker(2 * time.Minute),
		done:          make(chan bool),
	}
}
