package cli

import (
	"flag"
	"fmt"
	"github.com/pythoninja/go-redirect/internal/config"
	"github.com/pythoninja/go-redirect/internal/database"
	"github.com/pythoninja/go-redirect/internal/server"
	"github.com/pythoninja/go-redirect/internal/storage"
	"github.com/pythoninja/go-redirect/internal/vars"
	"log/slog"
	"os"
)

const (
	flagDisplayVersionHelp       = "Display app version and exit"
	flagServerAddrHelp           = "App address to listen"
	flagServerPortHelp           = "App server port to listen"
	flagEnvironmentHelp          = "Environment (development|production"
	flagDatabaseDsnHelp          = "Database DSN"
	flagDatabaseMaxOpenConnsHelp = "Postgres max open connections"
	flagDatabaseMaxIdleConnsHelp = "Postgres max idle connections"
	flagDatabaseMaxIdleTimeHelp  = "Postgres max time. Values example: 45s, 30m"
	flagEnableRateLimiterHelp    = "Enable global rate limiter"
	flagSetAPIKeyHelp            = "Set a custom API key or leave empty to generate a random"
)

//goland:noinspection GoUnhandledErrorResult
func Run() {
	var cfg config.Config

	flagVersion := flag.Bool("v", false, flagDisplayVersionHelp)
	flag.StringVar(&cfg.Env, "env", "development", flagEnvironmentHelp)
	flag.StringVar(&cfg.Addr, "addr", "0.0.0.0", flagServerAddrHelp)
	flag.IntVar(&cfg.Port, "port", 4000, flagServerPortHelp)
	flag.StringVar(&cfg.Database.Dsn, "db-dsn", os.Getenv("REDIRECT_DB_DSN"), flagDatabaseDsnHelp)
	flag.IntVar(&cfg.Database.MaxOpenConns, "db-max-open-conns", 25, flagDatabaseMaxOpenConnsHelp)
	flag.IntVar(&cfg.Database.MaxIdleConns, "db-max-idle-conns", 25, flagDatabaseMaxIdleConnsHelp)
	flag.StringVar(&cfg.Database.MaxIdleTime, "db-max-idle-time", "15m", flagDatabaseMaxIdleTimeHelp)
	flag.BoolVar(&cfg.EnableRateLimiter, "rate-limiter-enabled", true, flagEnableRateLimiterHelp)
	flag.StringVar(&cfg.APISecretKey, "api-key", "", flagSetAPIKeyHelp)
	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s v%s on %s\n", vars.Name, vars.Version, vars.Environment)
		os.Exit(0)
	}

	config.InitLogger() // todo: move logger into own package

	app := config.InitConfiguration(&cfg)

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
