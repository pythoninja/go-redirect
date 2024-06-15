# go-redirect

## Development

### Just commands

1. `just` - prints all commands.
2. `just switch <branch>` - fetches from origin and switches to requested branch.
3. `just pull-master` - switches to master branch and pulls from origin.
4. `just push-branch` - pushes to the origin with `--set-upstream origin HEAD` arg.
5. `just lint` - runs linters.
6. `just sql` - runs psql.
7. `just migrations-up` - run migration scripts.
8. `just migrations-down` - rollback migrations scripts.
9. `just prepare-dev` - runs Postgres container with Docker Compose.
10. `just run-dev` - runs `air`.
11. `just build-snapshot` - builds snapshot version.
12. `just changelog` - updates changelog using `cliff` and commits the changes.


### Docker compose

This project uses docker `compose.yml` profiles to set up required dependencies.

There are 3 profiles defined: `dev`, `simple-prod` and `migration`.

Use profile `dev` to set up just a Postgres database. The application will use default DSN to connect to it.

A short note on profile `simple-prod` - it will set up a containerized application with Postgres database. 

Running the project:
1. Run database migration scripts: `docker compose --profile migration run migration`
2. Run project in dev mode: `docker compose --profile dev up --detach`
3. Run project in simple-prod mode: `docker compose --profile simple-prod up --detach`

### Environment variables

```
REDIRECT_ENVIRONMENT = "development"
REDIRECT_API_SECRET_KEY = ""
REDIRECT_ENABLE_RATELIMITER = "false"

REDIRECT_DB_DSN = "postgres://postgres:password@172.40.0.10/redirect?sslmode=disable"
REDIRECT_DB_MAX_OPEN_CONNECTIONS = 25
REDIRECT_DB_MAX_IDLE_CONNECTIONS = 25
REDIRECT_DB_MAX_IDLE_TIME = "15m"

REDIRECT_LISTEN_ADDRESS = "0.0.0.0"
REDIRECT_LISTEN_PORT = 4000
```
