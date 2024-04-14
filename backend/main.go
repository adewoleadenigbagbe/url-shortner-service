package main

import (
	"os"

	"github.com/adewoleadenigbagbe/url-shortner-service/cmd"
)

func main() {
	shortner := cmd.NewUrlShortner()
	err := shortner.Start()
	if err != nil {
		os.Exit(1)
	}
}
