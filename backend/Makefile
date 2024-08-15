install:
	go mod tidy
	go install github.com/rubenv/sql-migrate/...@latest

generate:
	rm -rf pkg/reposqlc
	cd config && sqlc generate

migrate:
	sql-migrate up -config ./config/sql-migrate.yml

unmigrate:
	sql-migrate down -config ./config/sql-migrate.yml

run:
	sh scripts/run.sh

build:
	sh scripts/build.sh

dockerize:
	sh ./scripts/docker.sh
