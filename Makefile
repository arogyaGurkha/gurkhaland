FRONT_END_BINARY=frontApp
BROKER_SERVICE_BINARY=brokerApp
LOGGER_SERVICE_BINARY=loggerApp
AUTH_SERVICE_BINARY=authApp
MAIL_SERVICE_BINARY=mailApp
LISTENER_SERVICE_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker_service build_auth_service build_mail_service build_listener_service build_logger_service
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## build_broker_service: builds the broker service binary as a linux executable
build_broker_service:
	@echo "Building broker service binary..."
	cd broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${BROKER_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_logger_service: builds the logger service binary as a linux executable
build_logger_service:
	@echo "Building logger service binary..."
	cd logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${LOGGER_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_auth_service: builds the auth service binary as a linux executable
build_auth_service:
	@echo "Building auth service binary..."
	cd auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${AUTH_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_mail_service: builds the mail service binary as a linux executable
build_mail_service:
	@echo "Building mail service binary..."
	cd mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${MAIL_SERVICE_BINARY} ./cmd/api
	@echo "Done!"

## build_listener_service: builds the listener service binary as a linux executable
build_listener_service:
	@echo "Building listener service binary..."
	cd listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${LISTENER_SERVICE_BINARY} .
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd front-end && env CGO_ENABLED=0 go build -o ./bin/${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start_front: build_front
	@echo "Starting front end"
	cd front-end && ./bin/${FRONT_END_BINARY} &

## stop: stop the front end
stop_front:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"