.PHONY: run build test migrate-up migrate-down

APP_NAME=pizza-hub
DB_URL=postgres://pizza:pizza123@localhost:5432/pizzahub?sslmode=disable

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP_NAME) ./cmd/server

test:
	go test ./... -v

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

tidy:
	go mod tidy