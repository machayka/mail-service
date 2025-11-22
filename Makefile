run: build
	@./bin/app

dev:
	@$(shell go env GOPATH)/bin/CompileDaemon -command="./app" -build="go build -o app ./cmd/app"

build:
	@go build -o bin/app cmd/app/main.go
