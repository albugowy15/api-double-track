.PHONY: migrate_up migrate_down

migrate_up: 
	go run cmd/database/migrations/up/main.go

migrate_down: 
	go run cmd/database/migrations/down/main.go

