package cli

import (
	"flag"
	"fmt"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/database"
	"github.com/pythoninja/go-redirect/internal/server"
	"github.com/pythoninja/go-redirect/internal/storage"
	"github.com/pythoninja/go-redirect/internal/version"
	"log/slog"
	"os"
)

const (
	flagDisplayVersionHelp       = "Display app version and exit"
	flagEnvironmentHelp          = "Environment (development|production) (env: REDIRECT_ENVIRONMENT)"
	flagServerAddrHelp           = "App address to listen (env: REDIRECT_LISTEN_ADDRESS)"
	flagServerPortHelp           = "App server port to listen (env: REDIRECT_LISTEN_PORT)"
	flagDatabaseDsnHelp          = "Database DSN (env: REDIRECT_DB_DSN)"
	flagDatabaseMaxOpenConnsHelp = "Postgres max open connections (env: REDIRECT_DB_MAX_OPEN_CONNECTIONS)"
	flagDatabaseMaxIdleConnsHelp = "Postgres max idle connections (env: REDIRECT_DB_MAX_IDLE_CONNECTIONS)"
	flagDatabaseMaxIdleTimeHelp  = "Postgres max time. Values example: 45s, 30m (env: REDIRECT_DB_MAX_IDLE_TIME)"
	flagEnableRateLimiterHelp    = "Enable global rate limiter (env: REDIRECT_ENABLE_RATELIMITER)"
	flagSetAPIKeyHelp            = "Secret API key. Leave blank to generate new (env: REDIRECT_API_KEY)" //nolint:gosec
)

func Run() {
	var cfg config.Config

	flagVersion := flag.Bool("v", false, flagDisplayVersionHelp)
	flag.StringVar(&cfg.Env, "env",
		config.GetEnv("REDIRECT_ENVIRONMENT", "development"), flagEnvironmentHelp)
	flag.StringVar(&cfg.Addr, "addr",
		config.GetEnv("REDIRECT_LISTEN_ADDRESS", "0.0.0.0"), flagServerAddrHelp)
	flag.IntVar(&cfg.Port, "port",
		config.GetEnv("REDIRECT_LISTEN_PORT", 4000), flagServerPortHelp)
	flag.StringVar(&cfg.Database.Dsn, "db-dsn",
		config.GetEnv("REDIRECT_DB_DSN", ""), flagDatabaseDsnHelp)
	flag.IntVar(&cfg.Database.MaxOpenConns, "db-max-open-conns",
		config.GetEnv("REDIRECT_DB_MAX_OPEN_CONNECTIONS", 25), flagDatabaseMaxOpenConnsHelp)
	flag.IntVar(&cfg.Database.MaxIdleConns, "db-max-idle-conns",
		config.GetEnv("REDIRECT_DB_MAX_IDLE_CONNECTIONS", 25), flagDatabaseMaxIdleConnsHelp)
	flag.StringVar(&cfg.Database.MaxIdleTime, "db-max-idle-time",
		config.GetEnv("REDIRECT_DB_MAX_IDLE_TIME", "15m"), flagDatabaseMaxIdleTimeHelp)
	flag.BoolVar(&cfg.EnableRateLimiter, "rate-limiter-enabled",
		config.GetEnv("REDIRECT_ENABLE_RATELIMITER", true), flagEnableRateLimiterHelp)
	flag.StringVar(&cfg.APISecretKey, "api-key",
		config.GetEnv("REDIRECT_API_KEY", ""), flagSetAPIKeyHelp)
	flag.Parse()

	app := config.InitConfiguration(&cfg)

	if *flagVersion {
		buildInfo := version.GetBuildInfo()
		fmt.Printf(`%s
Version:     %s
Environment: %s
Commit:      %s
Build time:  %s
`,
			buildInfo.AppName,
			buildInfo.Version,
			app.Config.Env,
			buildInfo.VcsCommit,
			buildInfo.VcsTime,
		)

		os.Exit(0)
	}

	config.SetupLogger(app.Config.Env) // todo: move logger into own package

	db, err := database.NewConnectionPool(app)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	store := storage.New(db)

	err = server.Serve(app, store)
	if err != nil {
		slog.Error("failed to run server", slog.Any("error", err.Error()))

		return
	}
}
