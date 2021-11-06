include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/app: run the application
.PHONY: run/app
run/app:
	CHATAPP_DB_DSN=${CHATAPP_DB_DSN} go run ./cmd/web

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${CHATAPP_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration file for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations'
	migrate -path=./migrations -database=${CHATAPP_DB_DSN} up

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/app: build the application
.PHONY: build/app
build/app:
	@echo 'Building application...'
	@go1.16.5 build -o bin/chatapp ./cmd/web

## build/run: run the binary
.PHONY: build/run
build/run:
	@CHATAPP_DB_DSN=${CHATAPP_DB_DSN} ./bin/chatapp