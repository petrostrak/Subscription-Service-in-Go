package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func StartApp() {
	// connect to DB
	db := initDB()
	db.Ping()

	// create sessions

	// create channels

	// create waitGroup

	// setup the app config

	// setup mail

	// listen for web connections
}

func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to DB")
	}

	return conn
}

func connectToDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready")
		} else {
			log.Println("connnected to DB")
			return connection
		}

		if counts > 10 {
			return nil
		}

		log.Println("Backing off for a moment")
		counts++
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
