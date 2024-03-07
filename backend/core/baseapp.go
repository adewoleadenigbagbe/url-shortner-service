package core

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	services "github.com/adewoleadenigbagbe/url-shortner-service/service"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type BaseApp struct {
	Echo        *echo.Echo
	Db          *sql.DB
	Rdb         *redis.Client
	AuthService services.AuthService
	UrlService  services.UrlService
}

func ConfigureAppDependencies() (*BaseApp, error) {
	//connect to sqllite
	db, err := database.ConnectToSQLite()
	if err != nil {
		return nil, err
	}

	//connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ping, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ping)

	app := &BaseApp{
		Echo: echo.New(),
		Db:   db,
		Rdb:  redisClient,
		AuthService: services.AuthService{
			Db:  db,
			Rdb: redisClient,
		},
		UrlService: services.UrlService{
			Db: db,
		},
	}

	return app, nil
}
