package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adewoleadenigbagbe/url-shortner-service/key-generator-service/keyservice"
)

func main() {

	kg, err := keyservice.NewKeyGenerator()
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	//ctx := context.WithValue(context.Background(), "sig", quit)
	go func(c chan os.Signal) {
		kg.Run(c)
	}(quit)

	<-quit
	fmt.Println("Exiting the main")
}
