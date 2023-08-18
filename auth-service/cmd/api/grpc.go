package main

import (
	"context"
	"fmt"
	"github.com/arogyaGurkha/gurkhaland-proto/auth-service/auth"
	"github.com/arogyaGurkha/gurkhaland-proto/logger-service/logs"
	"github.com/arogyaGurkha/gurkhaland/auth-service/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Models data.Models
}

func (a *AuthServer) Authenticate(ctx context.Context, req *auth.AuthRequest) (*auth.AuthResponse, error) {
	input := req.GetAuthReq()

	// Check credentials
	u, err := a.Models.User.GetByEmail(input.Email)
	if err != nil {
		log.Println("Invalid credentials: Email")
		res := &auth.AuthResponse{
			Result: "failed",
			Data:   nil,
		}
		return res, err
	}

	valid, err := u.PasswordMatches(input.Password)
	if err != nil || !valid {
		log.Println("Invalid credentials: PW")
		res := &auth.AuthResponse{
			Result: "failed", // TODO: manage result messages
			Data:   nil,
		}
		return res, err
	}

	user := &auth.User{
		Id:        int32(u.ID), // TODO: change proto file to not do this
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Active:    int32(u.Active),
		CreatedAt: timestamppb.New(u.CreatedAt), // TODO: provides time strangely
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}

	err = a.logEventViaGRPC("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		log.Println("logging auth request failed", err)
		res := &auth.AuthResponse{
			Result: "failed", // TODO: manage result messages
			Data:   nil,
		}
		return res, err
	}

	log.Println("Logged in user:", u.Email)
	res := &auth.AuthResponse{
		Result: "success", // TODO: manage result messages
		Data:   user,
	}
	return res, nil
}

func (a *AuthServer) logEventViaGRPC(name, data string) error {
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: name,
			Data: data,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
