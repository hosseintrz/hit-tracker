package database

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"os"
)

func InitStore() (*sql.DB, error) {

	pgConnString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	var (
		db  *sql.DB
		err error
	)
	openDB := func() error {
		db, err = sql.Open("postgres", pgConnString)
		return err
	}
	log.Infof(`db config: %s`, pgConnString)
	log.Info("connecting to database")
	//err = backoff.Retry(openDB, backoff.NewExponentialBackOff())
	err = openDB()
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS visits (
    		user_ip STRING,
            url STRING,
            count INT,
            PRIMARY KEY (user_ip, url))`); err != nil {
		return nil, err
	}
	log.Info("connected to database")
	return db, nil
}
