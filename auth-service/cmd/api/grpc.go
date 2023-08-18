package main

import (
	"context"
	"github.com/arogyaGurkha/gurkhaland-proto/auth-service/auth"
	"github.com/arogyaGurkha/gurkhaland/auth-service/data"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
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

	// TODO: log authentication requests

	user := &auth.User{
		Id:        int32(u.ID), // TODO: change proto file to not do this
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Active:    int32(u.Active),
		CreatedAt: timestamppb.New(u.CreatedAt), // TODO: provides time strangely
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}

	log.Println("Logged in user:", u.Email)
	res := &auth.AuthResponse{
		Result: "success", // TODO: manage result messages
		Data:   user,
	}
	return res, nil
}
