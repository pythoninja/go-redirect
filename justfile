#!/usr/bin/env just --justfile

default:
    @just --list

switch branch:
    @echo 'Fetching and switching to {{branch}}'
    git fetch origin && git switch --create '{{branch}}'

pull-master:
    git switch master && git pull

push-branch:
    git push --set-upstream origin HEAD

prepare-dev:
    @docker compose --profile dev up --detach

run-dev:
    @air

sql:
    @psql $REDIRECT_DB_DSN

build-snapshot:
    @goreleaser build --snapshot --clean

[confirm('Run all migrations?')]
migrations-up:
    @migrate -path ./migrations -database=$REDIRECT_DB_DSN up

[confirm('Rollback all migrations?')]
migrations-down:
    @migrate -path ./migrations -database=$REDIRECT_DB_DSN down --all

lint:
	@echo 'Tidying and verifying module dependencies...'
	@go mod tidy
	@go mod verify
	@echo 'Formatting code...'
	@go fmt ./...
	@echo 'Vetting code...'
	@go vet ./...
	@echo 'Linting code with staticcheck...'
	@staticcheck ./...
	@echo 'Linting code with golangci-lint...'
	@golangci-lint run
	@echo 'Linting code with revive...'
	@revive -formatter friendly ./...
