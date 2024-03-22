.PHONY: migrate_up migrate_down seed server client dev test build start docs

migrate_up: 
	go run cmd/database/migrations/up/main.go

migrate_down: 
	go run cmd/database/migrations/down/main.go

seed:
	go run cmd/database/seeder/main.go

server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

dev:
	air

test:
	go test -v ./...

build:
	go build -o ./tmp/main cmd/server/main.go

start:
	go build -o ./tmp/main cmd/server/main.go && ./tmp/main

docs:
	swag init -d "./" -g "cmd/server/main.go"
