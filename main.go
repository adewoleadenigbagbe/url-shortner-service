package main

import (
	"fmt"
	"log"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
)

func main() {
	_, err := database.ConnectToSQLite()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}
