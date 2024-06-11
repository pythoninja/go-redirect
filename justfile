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

run:
    @air

sql:
    psql $REDIRECT_DB_DSN

[confirm('Run all migrations?')]
migrations-up:
    @migrate -path ./migrations -database=$REDIRECT_DB_DSN up

[confirm('Rollback all migrations?')]
migrations-down:
    @migrate -path ./migrations -database=$REDIRECT_DB_DSN down --all

audit:
	@echo 'Tidying and verifying module dependencies...'
	@go mod tidy
	@go mod verify
	@echo 'Formatting code...'
	@go fmt ./...
	@echo 'Vetting code...'
	@go vet ./...
	@echo 'Linting code with staticcheck...'
	@staticcheck ./...
