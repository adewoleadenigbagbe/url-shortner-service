package services

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	Db  *sql.DB
	Rdb *redis.Client
}

type UrlService struct {
	Db *sql.DB
}

type DomainService struct {
	Db *sql.DB
}
