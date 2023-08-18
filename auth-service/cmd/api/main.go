package main

import (
	"database/sql"
	"fmt"
	"github.com/arogyaGurkha/gurkhaland-proto/auth-service/auth"
	"github.com/arogyaGurkha/gurkhaland/auth-service/data"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort  = "80"
	gRpcPort = "50001"
)

var retryCounts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting Postgres connection")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Failed connecting to Postgres")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	go app.gRPCListener()
	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Println("Starting authentication service on port", webPort)

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

func (app *Config) gRPCListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", gRpcPort)) // TODO: hardcoded tcp port
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	auth.RegisterAuthServiceServer(s, &AuthServer{Models: app.Models}) // TODO: Simplify GRPC Server call

	log.Printf("gRPC server started on port %v", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
