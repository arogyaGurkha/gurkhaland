module github.com/arogyaGurkha/gurkhaland/broker-service

go 1.20

require (
	github.com/arogyaGurkha/gurkhaland-proto/auth-service v0.0.0-20230818064544-44701dc66c93
	github.com/arogyaGurkha/gurkhaland-proto/logger-service v0.0.0-20230818063500-ab675c66077b
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/cors v1.2.1
	github.com/rabbitmq/amqp091-go v1.8.1
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/arogyaGurkha/gurkhaland-proto/mail-service v0.0.0-20230818072117-5db41e155737 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
)
