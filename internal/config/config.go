package config

import (
	"log/slog"
	"os"
)

type Application struct {
	Config Config
}

type Config struct {
	Port     int
	Env      string
	Database struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

func InitConfiguration(cfg *Config) *Application {
	return &Application{
		Config: Config{
			Env:      cfg.Env,
			Port:     cfg.Port,
			Database: cfg.Database,
		},
	}
}

func InitLogger() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	log.Info("initialize default logger")
	slog.SetDefault(log)
}
