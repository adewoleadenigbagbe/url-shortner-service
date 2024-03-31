package core

import (
	"context"
	"database/sql"
	"os"
	"time"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	auth "github.com/adewoleadenigbagbe/url-shortner-service/services/auth"
	domain "github.com/adewoleadenigbagbe/url-shortner-service/services/domain"
	link "github.com/adewoleadenigbagbe/url-shortner-service/services/shortlinks"
	user "github.com/adewoleadenigbagbe/url-shortner-service/services/user"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	dbFile = "urlshortnerDB.db"
)

type BaseApp struct {
	Echo          *echo.Echo
	Db            *sql.DB
	Rdb           *redis.Client
	AuthService   auth.AuthService
	UrlService    link.UrlService
	DomainService domain.DomainService
	UserService   user.UserService
}

func ConfigureAppDependencies() (*BaseApp, error) {
	//connect to sqllite
	db, err := database.ConnectToSQLite(dbFile)
	if err != nil {
		return nil, err
	}

	//connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	app := &BaseApp{
		Echo: echo.New(),
		Db:   db,
		Rdb:  redisClient,
		AuthService: auth.AuthService{
			Db:  db,
			Rdb: redisClient,
		},
		UrlService: link.UrlService{
			Db: db,
		},
		DomainService: domain.DomainService{
			Db: db,
		},
		UserService: user.UserService{
			Db: db,
		},
	}

	return app, nil
}
