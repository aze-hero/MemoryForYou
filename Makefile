.PHONY: dev dev-web dev-api db-up db-down build clean

dev:
	pnpm dev

dev-web:
	pnpm --filter=@timecapsule/web dev

dev-api:
	cd apps/api && go run ./cmd/server/main.go

db-up:
	docker run -d --name timecapsule-db \
		-e POSTGRES_USER=timecapsule \
		-e POSTGRES_PASSWORD=timecapsule \
		-e POSTGRES_DB=timecapsule \
		-p 5432:5432 \
		docker.1ms.run/library/postgres:16-alpine 2>/dev/null || echo "DB already running"

db-down:
	docker rm -f timecapsule-db 2>/dev/null || echo "DB not running"

db-migrate:
	psql "postgresql://timecapsule:timecapsule@localhost:5432/timecapsule?sslmode=disable" -f apps/api/migrations/001_init.up.sql

build:
	pnpm build
	cd apps/api && go build -o bin/api ./cmd/server/main.go

clean:
	pnpm clean
	rm -rf apps/api/bin
