package helpers

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

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

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func IsValidUrl(query string) bool {
	_, err := url.ParseRequestURI(query)
	return err == nil
}
