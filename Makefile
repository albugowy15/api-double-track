.PHONY: migrate_up migrate_down seed server client dev build start

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

build:
	go build -o ./tmp/main cmd/server/main.go

start:
	go build -o ./tmp/main cmd/server/main.go && ./tmp/main
