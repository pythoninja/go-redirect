# go-redirect

## Development

Just commands:
1. `just switch <branch>` - fetches from origin and switches to requested branch.
2. `just master-pull` - switches to master branch and pulls from origin.
3. `just branch-push` - pushes to the origin with `--set-upstream origin HEAD` arg.

Connect to postgresql:

`psql $REDIRECT_DB_DSN` - will connect to the PostgreSQL instance running with `compose.yml`.