.PHONY: migrate_up migrate_down seed server client dev test build start docs

migrate_up: 
	go run cmd/database/migrations/up/main.go

migrate_down: 
	go run cmd/database/migrations/down/main.go

seed:
	go run cmd/database/seeder/main.go

run:
	go run cmd/api/main.go

dev:
	air

test:
	go test -v ./...

build:
	go build -o ./tmp/main cmd/api/main.go

start:
	go build -o ./tmp/main cmd/api/main.go && ./tmp/main

docs:
	swag init -d "./" -g "cmd/api/main.go"
