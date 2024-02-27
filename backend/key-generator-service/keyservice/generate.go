package keyservice

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
)

const (
	characterLimit = 6
	expirySpan     = 1
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
	db, err := database.ConnectToSQLite()
	if err != nil {
		log.Fatal(err)
	}

	return &KeyGenerator{
		db:            db,
		scheduleTimer: time.NewTicker(2 * time.Minute),
		done:          make(chan bool),
	}
}