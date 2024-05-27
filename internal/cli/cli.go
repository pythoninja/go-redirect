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
	flagServerPortHelp           = "App server port to listen"
	flagEnvironmentHelp          = "Environment (development|production"
	flagDatabaseDsnHelp          = "Database DSN"
	flagDatabaseMaxOpenConnsHelp = "Postgres max open connections"
	flagDatabaseMaxIdleConnsHelp = "Postgres max idle connections"
	flagDatabaseMaxIdleTimeHelp  = "Postgres max time. Values example: 45s, 30m"
)

//goland:noinspection GoUnhandledErrorResult
func Run() {
	var cfg config.Config

	flagVersion := flag.Bool("v", false, flagDisplayVersionHelp)
	flag.IntVar(&cfg.Port, "port", 4000, flagServerPortHelp)
	flag.StringVar(&cfg.Env, "env", "development", flagEnvironmentHelp)
	flag.StringVar(&cfg.Database.Dsn, "db-dsn", os.Getenv("REDIRECT_DB_DSN"), flagDatabaseDsnHelp)
	flag.IntVar(&cfg.Database.MaxOpenConns, "db-max-open-conns", 25, flagDatabaseMaxOpenConnsHelp)
	flag.IntVar(&cfg.Database.MaxIdleConns, "db-max-idle-conns", 25, flagDatabaseMaxIdleConnsHelp)
	flag.StringVar(&cfg.Database.MaxIdleTime, "db-max-idle-time", "15m", flagDatabaseMaxIdleTimeHelp)
	flag.Parse()

	if *flagVersion {
		fmt.Printf("%s v%s on %s\n", vars.Name, vars.Version, vars.Environment)
		os.Exit(0)
	}

	db, err := database.NewConnectionPool(&cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := config.InitConfiguration(&cfg)
	store := storage.New(db)
	config.InitLogger() // todo: move logger into own package

	err = server.Serve(app, store)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
