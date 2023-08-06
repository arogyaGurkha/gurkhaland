package main

import (
	"context"
	"github.com/arogyaGurkha/gurkhaland/logger-service/data"
	"log"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

// LogInfo logs RPCPayload information to MongoDB collection, sets resp string pointer to success message
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs") // TODO: hardcoded DB name
	result, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo with RPC", err)
		return err
	}
	log.Println("Logged item", result)

	*resp = "Processed payload via RPC:" + payload.Name

	return nil
}
