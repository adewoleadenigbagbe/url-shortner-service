package jwtauth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// generate JWT token
func GenerateJWT(user models.SigInUserRequest) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
		"iat":   time.Now().Unix(),
		"eat":   time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})

	privateKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))
	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(context echo.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// // check token validity
func getToken(context echo.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		privateKey := []byte(os.Getenv("JWT_PRIVATE_KEY"))
		return privateKey, nil
	})
	return token, err
}

// // extract token from request Authorization header
func getTokenFromRequest(context echo.Context) string {
	bearerToken := context.Request().Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}

	return ""
}
