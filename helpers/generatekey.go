package helpers

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateApiKey(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	//taking just 20 characters from the hash, as the first 15 are mostly duplicate from the hash
	return base64.StdEncoding.EncodeToString(hash)[15:35]
}
