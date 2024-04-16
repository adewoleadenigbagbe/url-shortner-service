package linkservice

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	//ctx := context.WithValue(context.Background(), "sig", quit)

	sg := NewShortlinkGenerator()
	go func(c chan os.Signal) {
		sg.GenerateLink(c)
	}(quit)

	<-quit
	fmt.Println("Exiting the main")
}
