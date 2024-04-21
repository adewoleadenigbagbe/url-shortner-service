package core

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	auth "github.com/adewoleadenigbagbe/url-shortner-service/services/auth"
	domain "github.com/adewoleadenigbagbe/url-shortner-service/services/domain"
	enum "github.com/adewoleadenigbagbe/url-shortner-service/services/enums"
	export "github.com/adewoleadenigbagbe/url-shortner-service/services/export"
	link "github.com/adewoleadenigbagbe/url-shortner-service/services/shortlinks"
	statistic "github.com/adewoleadenigbagbe/url-shortner-service/services/statistics"
	tag "github.com/adewoleadenigbagbe/url-shortner-service/services/tags"
	teams "github.com/adewoleadenigbagbe/url-shortner-service/services/teams"
	user "github.com/adewoleadenigbagbe/url-shortner-service/services/user"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	dbFile = "data/urlshortnerDB.db"
)

type BaseApp struct {
	Echo              *echo.Echo
	Db                *sql.DB
	Rdb               *redis.Client
	AuthService       auth.AuthService
	UrlService        link.UrlService
	DomainService     domain.DomainService
	UserService       user.UserService
	TeamService       teams.TeamService
	TagService        tag.TagService
	ExportService     export.ExportService
	StatisticsService statistic.StatisticsService
	EnumService       enum.EnumService
}

func ConfigureAppDependencies() (*BaseApp, error) {
	//connect to sqlite
	db, err := database.ConnectToSQLite(dbFile)
	if err != nil {
		fmt.Println("here aaaa")
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
		TeamService: teams.TeamService{
			Db: db,
		},
		TagService: tag.TagService{
			Db: db,
		},
		ExportService: export.ExportService{
			Db: db,
		},
		StatisticsService: statistic.StatisticsService{
			Db: db,
		},
		EnumService: enum.EnumService{},
	}

	return app, nil
}
