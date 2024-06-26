include .env

export DATABASE_URL ?= postgres://$(SQL_USERNAME):$(SQL_PASSWORD)@$(SQL_HOST):$(SQL_PORT)/$(SQL_DATABASE)?sslmode=$(SQL_SSL)

bin:
	@mkdir -p bin

setup-tools: bin
	# @curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
ifeq ($(shell uname), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else ifeq ($(shell uname), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate bin/
else
	@echo "Your OS is not supported."
endif

migration-create:
	bin/migrate create -ext sql -dir migrations -seq $(name)

migration-up:
	bin/migrate -path migrations -database "${DATABASE_URL}" up

migration-down:
	bin/migrate -path migrations -database "${DATABASE_URL}" down $(n)

seed-create:
	bin/migrate create -ext sql -dir migrations/seeds -seq $(name)

seed-up:
	bin/migrate -path migrations/seeds -database "${DATABASE_URL}&x-migrations-table=seed_migrations" up

seed-down:
	bin/migrate -path migrations/seeds -database "${DATABASE_URL}&x-migrations-table=seed_migrations" down $(n)

run-dev:
	bin/air

run:
	./main

build:
	go build ./main.go

test:
	go test -v -cover ./...

download-file:
	@echo "Downloading file from Google Drive"
	curl -L -o datasets.parquet "https://drive.usercontent.google.com/download?id=1QLBGFOoKw_3-iM58q4unWfwHmPqfnrYr&export=download&authuser=0&confirm=t&uuid=0e0a5b66-1b23-4d93-a9f8-4814b82a26ae&at=APZUnTWQH6Py1QF0cizaF7hzYyIJ%3A1719379161030"


