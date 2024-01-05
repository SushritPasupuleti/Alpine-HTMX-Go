include .env

run:
	@echo "Running server..."
	@go run main.go

build:
	@echo "Building server..."
	@GOOS=linux GOARCH=amd64 go build -o bin/server main.go 
	@GOOS=windows GOARCH=amd64 go build -o bin/server.exe main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/server.darwin main.go
	@echo "Done building for linux, windows and mac (x64 only)"

test:
	@echo "Running tests..."
	@go test -v ./...

migrate_up:
	@echo "Migrating up..."
	sqlx migrate run --database-url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

packages_install:
	@echo "Installing packages..."
	@go mod tidy

packages_update:
	@echo "Updating packages..."
	@go get -u

dev:
	@air
