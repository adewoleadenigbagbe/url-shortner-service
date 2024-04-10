package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"math"
	"math/rand"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func encodeToString(text string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(hash)
}

func GenerateApiKey(text string) string {
	key := encodeToString(text)
	//taking just 20 characters from the hash, as the first 15 are mostly duplicate from the hash
	return key[15:35]
}

func GenerateShortLink(text string) string {
	algorithm := sha256.New()
	algorithm.Write([]byte(text))
	hash := algorithm.Sum(nil)
	short := base64.StdEncoding.EncodeToString(hash)

	return short[:8]
}

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func GeneratePassword(text string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	return string(hash)
}

func IsValidUrl(query string) bool {
	_, err := url.ParseRequestURI(query)
	return err == nil
}

func StartOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, nil)
}

func EndOfDay(date time.Time) time.Time {
	nsec := int(math.Pow10(9)) - 1
	return time.Date(date.Year(), date.Month(), date.Day(), 11, 59, 59, nsec, nil)
}
