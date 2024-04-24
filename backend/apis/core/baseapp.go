package core

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	auth "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/auth"
	domain "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/domain"
	enum "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/enums"
	export "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/export"
	plan "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/payplan"
	link "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/shortlinks"
	statistic "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/statistics"
	tag "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/tags"
	teams "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/teams"
	user "github.com/adewoleadenigbagbe/url-shortner-service/apis/services/user"

	database "github.com/adewoleadenigbagbe/url-shortner-service/db"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const (
	dbFile = "./data/urlshortnerDB.db"
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
	PlanService       plan.PlanService
}

func ConfigureAppDependencies() (*BaseApp, error) {
	path, err := filepath.Abs(dbFile)
	if err != nil {
		return nil, err
	}

	//connect to sqlite
	db, err := database.ConnectToSQLite(path)
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
		PlanService: plan.PlanService{
			Db: db,
		},
	}

	return app, nil
}
