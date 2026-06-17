.PHONY: run migrate-up migrate-down test clean dev dev-stop docker-build docker-up docker-down docker-logs
run:
	go run cmd/server/main.go
dev:
	docker-compose up -d
dev-stop:
	docker-compose down
migrate-up:
	migrate -path ./migrations -database "postgres://admin:secret@localhost:5432/expensedb?sslmode=disable" up
migrate-down:
	migrate -path ./migrations -database "postgres://admin:secret@localhost:5432/expensedb?sslmode=disable" down
test:
	go test ./...
clean:
	docker-compose down -v
docker-build:
	docker compose build
docker-up:
	docker compose up --build -d
docker-down:
	docker compose down
docker-logs:
	docker compose logs -f app
