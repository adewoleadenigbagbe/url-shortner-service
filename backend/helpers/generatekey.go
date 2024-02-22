package helpers

import (
	"encoding/base64"
	"log"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

func encodeToString(text string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	//taking just 20 characters from the hash, as the first 15 are mostly duplicate from the hash
	return base64.StdEncoding.EncodeToString(hash)[15:35]
}

func GenerateApiKey(text string) string {
	key := encodeToString(text)
	//taking just 20 characters from the hash, as the first 15 are mostly duplicate from the hash
	return key[15:35]
}

func GenerateShortUrl(text string) string {
	short := encodeToString(text)

	//taking just 6 characters from the hash, as the first 15 are mostly duplicate from the hash
	return short[15:21]
}

func IsValidUrl(query string) bool {
	_, err := url.ParseRequestURI(query)
	return err == nil
}
