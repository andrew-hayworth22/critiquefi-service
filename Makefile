build:
	docker build \
	  --build-arg VERSION=1.0.0 \
	  --build-arg COMMIT=$(git rev-parse --short HEAD) \
	  -t critiquefi-service:latest \
	  -t critiquefi-service:1.0.0 .

run:
	docker stop critiquefi-service && docker rm critiquefi-service && docker run -d --env-file .env -p 8080:8080 --name critiquefi-service critiquefi-service

genkey:
	go run ./cmd/genkey/main.go

migrate-up:
	go run ./cmd/migrate -action=up

migrate-down:
	go run ./cmd/migrate -action=down

migrate-fresh:
	go run ./cmd/migrate -action=fresh