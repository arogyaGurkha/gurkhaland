package main

import (
	"context"
	"fmt"
	"github.com/arogyaGurkha/gurkhaland/logger-service/data"
	logs2 "github.com/arogyaGurkha/gurkhaland/logger-service/proto/logs"
	"google.golang.org/grpc"
	"log"
	"net"
)

type LogServer struct {
	logs2.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs2.LogRequest) (*logs2.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs2.LogResponse{Result: "failed"}
		return res, err
	}

	// return response
	res := &logs2.LogResponse{Result: "logged"}
	return res, nil
}

func (app *Config) gRPCListener() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", gRpcPort)) // TODO: hardcoded tcp port
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	logs2.RegisterLogServiceServer(s, &LogServer{Models: app.Models}) // TODO: Simplify GRPC Server call

	log.Printf("gRPC server started on port %v", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
