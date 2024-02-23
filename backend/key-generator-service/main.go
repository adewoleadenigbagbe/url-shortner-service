package main

import (
	"context"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/key-generator-service/keyservice"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kg := keyservice.NewKeyGenerator()
	go func() {
		kg.Run(ctx)
	}()

	// signal.Notify(kg.done, os.Interrupt, syscall.SIGTERM)
	// <-kg.done
	// fmt.Println("Exiting key generator service ....")
}
