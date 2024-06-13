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

### Docker compose

This project uses docker `compose.yml` profiles to set up required dependencies.

There are 3 profiles defined: `dev`, `simple-prod` and `migration`.

Use profile `dev` will set up just a Postgres database. The application will use default DSN to connect to it.

A short note on profile `simple-prod` - it will set up a containerized application with Postgres database. 

Running the project:
1. Run database migration scripts: `docker compose --profile migration run migration`
2. Run project in dev mode: `docker compose --profile dev up --detach`
3. Run project in simple-prod mode: `docker compose --profile simple-prod up --detach`
