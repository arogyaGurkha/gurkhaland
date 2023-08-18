package main

import (
	"fmt"
	"github.com/arogyaGurkha/gurkhaland/mail-service/proto/mail"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const (
	webPort  = "80"
	gRPCPort = "50001"
)

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port: ", webPort)
	go app.gRPCListener()
	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}

func (app *Config) gRPCListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", gRPCPort)) // TODO: hardcoded tcp port
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	mail.RegisterMailServiceServer(s, &MailServer{Mailer: app.Mailer})
	log.Printf("gRPC server started on port %v", gRPCPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
