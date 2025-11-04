run:
	go run ./cmd/api

genkey:
	go run ./cmd/genkey/main.go

migrate-up:
	go run ./cmd/migrate -action=up

migrate-down:
	go run ./cmd/migrate -action=down

migrate-fresh:
	go run ./cmd/migrate -action=fresh