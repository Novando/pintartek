install:
	go mod tidy
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	cd config && sqlc generate

generate:
	rm -rf pkg/reposqlc
	cd config && sqlc generate

migrate:
	sql-migrate up -config ./configs/sql-migrate.yml

unmigrate:
	sql-migrate down -config ./configs/sql-migrate.yml

run:
	sh scripts/run.sh

build:
	sh scripts/build.sh
dockerize:
	IMAGE="purify" sh ./scripts/docker.sh
