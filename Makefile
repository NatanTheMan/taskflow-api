include .env
export

run: cmd/api/main.go
	go run cmd/api/main.go

migrate-up:
	migrate -path migrations -database "$$DATABASE_URL" up
migrate-down:
	migrate -path migrations -database "$$DATABASE_URL" down 1
