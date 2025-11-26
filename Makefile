include .env
export

run: build
	@./bin/app

dev:
	@$(shell go env GOPATH)/bin/CompileDaemon -command="./app" -build="go build -o app ./cmd/app"

build:
	@go build -o bin/app cmd/app/main.go

# wywołanie: make migrate m=init.sql
migrate:
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f migrations/$(m) 
