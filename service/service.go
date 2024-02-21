package services

import "database/sql"

type AuthService struct {
	Db *sql.DB
}

type UrlService struct {
	Db *sql.DB
}
