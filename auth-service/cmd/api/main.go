package main

import (
	"database/sql"
	"fmt"
	"github.com/arogyaGurkha/gurkhaland/auth-service/data"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var retryCounts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("Starting authentication service \n")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Failed connecting to Postgres")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Printf("Postgres not yet ready...\n")
			retryCounts++
		} else {
			log.Printf("Connected to Postgres \n")
			return connection
		}

		if retryCounts > 10 {
			log.Println(err)
			return nil
		}

		log.Printf("Backing off...\n")
		time.Sleep(3 * time.Second)
		continue
	}
}
