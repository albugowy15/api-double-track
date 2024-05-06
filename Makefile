migrate_up: 
	@go run cmd/database/migrations/up/main.go

migrate_down: 
	@go run cmd/database/migrations/down/main.go

seed:
	@go run cmd/database/seeder/main.go

setup_db:
	@docker-compose up -d db && make migrate_down && make migrate_up && make dev

run:
	@go run cmd/api/main.go

dev:
	@docker-compose up -d db && air

test:
	@go test -v ./...

build:
	@go build -o ./tmp/main cmd/api/main.go

start:
	@go build -o ./tmp/main cmd/api/main.go && ./tmp/main

doc:
	@swag init -d "./" -g "cmd/api/main.go"
