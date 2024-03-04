.PHONY: migrate_up migrate_down seed

migrate_up: 
	go run cmd/database/migrations/up/main.go

migrate_down: 
	go run cmd/database/migrations/down/main.go

seed:
	go run cmd/database/seeder/main.go
