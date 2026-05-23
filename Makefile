.PHONY: dev dev-web dev-api db-up db-down build clean

dev:
	pnpm dev

dev-web:
	pnpm --filter=@timecapsule/web dev

dev-api:
	cd apps/api && go run ./cmd/server/main.go

db-up:
	docker compose -f docker/docker-compose.yml up -d

db-down:
	docker compose -f docker/docker-compose.yml down

db-migrate:
	psql "postgresql://timecapsule:timecapsule@localhost:5432/timecapsule?sslmode=disable" -f apps/api/migrations/001_init.up.sql

build:
	pnpm build
	cd apps/api && go build -o bin/api ./cmd/server/main.go

clean:
	pnpm clean
	rm -rf apps/api/bin
