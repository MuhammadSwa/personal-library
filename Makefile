include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage: make <target>'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/web
run/web:
	tailwindcss -i ./web/assets/css/main.css -o ./web/static/css/tailwind.css --watch &
	gow -c -e=go,mod,tmpl run ./cmd/web

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## db/migrations/new name=$1: create a new database migration file
.PHONY: db/migrations/new
migrations/new:
	@echo 'Creating migration file for ${DB_NAME}...'
	migrate create -seq -ext=.sql -dir=./internal/sql/schema ${DB_NAME}

## db/migrations/up: apply all up database migrations
.PHONY: migrations/up
migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./internal/sql/schema -database ${DB_DSN} up

## db/migrations/up: apply all down database migrations
.PHONY: migrations/down
migrations/down: confirm
	@echo 'Running up migrations...'
	migrate -path ./internal/sql/schema -database ${DB_DSN} down


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## TODO: make github actions for this
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

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'
git_description = $(shell git describe --always --dirty --tags --long)
## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/web
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/web
