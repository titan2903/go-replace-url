package database

import (
	"fmt"
	"log"
	"replace-url-gin/config"

	"github.com/jmoiron/sqlx"
)

func Postgres() *sqlx.DB {
	config := config.Get()

	connection := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=disable", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DB)
	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		log.Fatal("Error connection")
	}

	db.SetConnMaxIdleTime(config.DBConnectionIdle)
	db.SetConnMaxLifetime(config.DBConnectionLifetime)
	db.SetMaxIdleConns(config.DBMaxIdle)
	db.SetMaxOpenConns(config.DBMaxOpen)

	return db
}
