package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	billing "github.com/adewoleadenigbagbe/url-shortner-service/billing-service/service.go"
)

func main() {
	service, err := billing.NewBillingService()
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func(c chan os.Signal) {
		service.Run(c)
	}(quit)

	<-quit
	fmt.Println("Exiting the main")
}
