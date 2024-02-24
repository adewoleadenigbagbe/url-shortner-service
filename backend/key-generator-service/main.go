package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/adewoleadenigbagbe/url-shortner-service/key-generator-service/keyservice"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	//ctx := context.WithValue(context.Background(), "sig", quit)

	kg := keyservice.NewKeyGenerator()
	go func(c chan os.Signal) {
		fmt.Println("Starting the key generation service..")
		kg.Run(c)
	}(quit)

	<-quit
	fmt.Println("Exiting the main")
}
