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
	flagEnvironmentHelp          = "Environment (development|production)"
	flagServerAddrHelp           = "App address to listen"
	flagServerPortHelp           = "App server port to listen"
	flagDatabaseDsnHelp          = "Database DSN"
	flagDatabaseMaxOpenConnsHelp = "Postgres max open connections"
	flagDatabaseMaxIdleConnsHelp = "Postgres max idle connections"
	flagDatabaseMaxIdleTimeHelp  = "Postgres max time. Values example: 45s, 30m"
	flagEnableRateLimiterHelp    = "Enable global rate limiter"
	flagSetAPIKeyHelp            = "Set a custom API key or leave empty to generate a random" //nolint:gosec
)

func Run() {
	var cfg config.Config

	flagVersion := flag.Bool("v", false, flagDisplayVersionHelp)
	flag.StringVar(&cfg.Env, "env", os.Getenv("REDIRECT_ENVIRONMENT"), flagEnvironmentHelp)
	flag.StringVar(&cfg.Addr, "addr", "0.0.0.0", flagServerAddrHelp)
	flag.IntVar(&cfg.Port, "port", 4000, flagServerPortHelp)
	flag.StringVar(&cfg.Database.Dsn, "db-dsn", os.Getenv("REDIRECT_DB_DSN"), flagDatabaseDsnHelp)
	flag.IntVar(&cfg.Database.MaxOpenConns, "db-max-open-conns", 25, flagDatabaseMaxOpenConnsHelp)
	flag.IntVar(&cfg.Database.MaxIdleConns, "db-max-idle-conns", 25, flagDatabaseMaxIdleConnsHelp)
	flag.StringVar(&cfg.Database.MaxIdleTime, "db-max-idle-time", "15m", flagDatabaseMaxIdleTimeHelp)
	flag.BoolVar(&cfg.EnableRateLimiter, "rate-limiter-enabled", true, flagEnableRateLimiterHelp)
	flag.StringVar(&cfg.APISecretKey, "api-key", "", flagSetAPIKeyHelp)
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

	config.InitLogger() // todo: move logger into own package

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
