# go-redirect

## Development

Just commands:
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

